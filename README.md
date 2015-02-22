# gente
A golang websocket messaging abstraction

Patterned after [sente](https://github.com/ptaoussanis/sente), but not api compatible (as far as I know.)
Hopefully a relatively lightweight, understandable abstraction over websockets. 

Dependent upon [logrus](https://github.com/Sirupsen/logrus) and [gorilla/websocket](https://github.com/gorilla/websocket)


Totally not anywhere near ready, I'm still building out the actual functionality. 
If you're as crazy as I am, use something like:
  
    package main
    
    import (
    	"github.com/Sirupsen/logrus"
    	"github.com/byatlas/gente"
    	"net/http"
    )
    
    var log *logrus.Logger
    
    func main() {
    	log := logrus.StandardLogger()
    	routerBuilder := gente.RouterBuilder{}
    
    	msgPipe := gente.JsonCallbackPipeline{
    		Log:    log,
    		Router: routerBuilder.Finalize(),
    	}
    
    	http.Handle("/ws", gente.NewConnection(&msgPipe, *log))
    
    	log.Fatal(http.ListenAndServe(":8080", nil))
    }

No javascript lib as of yet, but it'll hopefully involve some code generation.
