package watcher

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/datawire/teleproxy/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type empty struct{}

type listWatchAdapter struct {
	resource dynamic.NamespaceableResourceInterface
}

func (lw listWatchAdapter) List(options v1.ListOptions) (runtime.Object, error) {
	return lw.resource.List(options)
}

func (lw listWatchAdapter) Watch(options v1.ListOptions) (watch.Interface, error) {
	return lw.resource.Watch(options)
}

type Watcher struct {
	config       *rest.Config
	resources    []*v1.APIResourceList
	stores       map[string]cache.Store
	stop         chan empty
	stopChans    []chan struct{}
	stoppedChans []chan empty
}

// NewWatcher returns a Kubernetes Watcher for the specified cluster
func NewWatcher(kubeinfo *k8s.KubeInfo) *Watcher {
	config, err := kubeinfo.GetRestConfig()
	if err != nil {
		log.Fatalln("Failed to get REST config:", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := disco.ServerResources()
	if err != nil {
		log.Fatal(err)
	}

	w := &Watcher{
		config:    config,
		resources: resources,
		stores:    make(map[string]cache.Store),
		stop:      make(chan empty),
	}

	go func() {
		<-w.stop
		for _, c := range w.stopChans {
			close(c)
		}
		for {
			<-w.stop
		}
	}()

	return w
}

func (w *Watcher) resolve(resource string) (string, string, v1.APIResource) {
	resource = strings.ToLower(resource)
	if resource == "" {
		return "", "", v1.APIResource{}
	}
	for _, rl := range w.resources {
		for _, r := range rl.APIResources {
			candidates := []string{
				r.Name,
				r.Kind,
				r.SingularName,
			}
			candidates = append(candidates, r.ShortNames...)

			for _, c := range candidates {
				if resource == strings.ToLower(c) {
					var group string
					var version string
					parts := strings.Split(rl.GroupVersion, "/")
					switch len(parts) {
					case 1:
						group = ""
						version = parts[0]
					case 2:
						group = parts[0]
						version = parts[1]
					default:
						panic("unrecognized GroupVersion")
					}
					return group, version, r
				}
			}
		}
	}
	return "", "", v1.APIResource{}
}

func (w *Watcher) Canonical(name string) string {
	parts := strings.Split(name, "/")

	var kind string
	switch len(parts) {
	case 1:
		kind = parts[0]
		name = ""
	case 2:
		kind = parts[0]
		name = parts[1]
	default:
		return ""
	}

	_, _, res := w.resolve(kind)
	kind = strings.ToLower(res.Kind)

	if name == "" {
		return kind
	}

	if res.Namespaced {
		var namespace string

		parts = strings.Split(name, ".")
		switch len(parts) {
		case 1:
			namespace = "default"
		case 2:
			name = parts[0]
			namespace = parts[1]
		default:
			return ""
		}

		return fmt.Sprintf("%s/%s.%s", kind, name, namespace)
	} else {
		return fmt.Sprintf("%s/%s", kind, name)
	}
}

func (w *Watcher) Watch(resources string, listener func(*Watcher)) error {
	group, version, res := w.resolve(resources)
	if version == "" {
		return fmt.Errorf("unknown resource: '%s'", resources)
	}

	dyn, err := dynamic.NewForConfig(w.config)
	if err != nil {
		log.Fatal(err)
	}

	resource := dyn.Resource(schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: res.Name,
	})

	store, controller := cache.NewInformer(
		listWatchAdapter{resource},
		nil,
		5*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				listener(w)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				listener(w)
			},
			DeleteFunc: func(obj interface{}) {
				listener(w)
			},
		},
	)

	w.stores[w.Canonical(res.Kind)] = store

	stopChan := make(chan struct{})
	stoppedChan := make(chan empty)
	w.stoppedChans = append(w.stoppedChans, stoppedChan)
	go func() {
		controller.Run(stopChan)
		close(stoppedChan)
	}()
	w.stopChans = append(w.stopChans, stopChan)

	return nil
}

func (w *Watcher) List(kind string) []Resource {
	kind = w.Canonical(kind)
	store, ok := w.stores[kind]
	if ok {
		objs := store.List()
		result := make([]Resource, len(objs))
		for idx, obj := range objs {
			result[idx] = obj.(*unstructured.Unstructured).UnstructuredContent()
		}
		return result
	} else {
		return nil
	}
}

func (w *Watcher) Get(kind, qname string) Resource {
	resources := w.List(kind)
	for _, res := range resources {
		if strings.ToLower(res.QName()) == strings.ToLower(qname) {
			return res
		}
	}
	return Resource{}
}

func (w *Watcher) Exists(kind, qname string) bool {
	return w.Get(kind, qname).Name() != ""
}

func (w *Watcher) Stop() {
	w.stop <- empty{}
}

func (w *Watcher) Wait() {
	for _, c := range w.stoppedChans {
		<-c
	}
}
