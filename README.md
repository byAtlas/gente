# gente
A golang websocket messaging abstraction

Patterned after [sente](https://github.com/ptaoussanis/sente), but not api compatible (as far as I know.)
Hopefully a relatively lightweight, understandable abstraction over websockets. 

Dependent upon [logrus](https://github.com/Sirupsen/logrus) and [gorilla/websocket](https://github.com/gorilla/websocket)


Totally not anywhere near ready, I'm still building out the actual functionality. 
If you're as crazy as I am, use something like:
  
    func main() {
    	log := logrus.StandardLogger()
    	routerBuilder := gente.RouterBuilder{}
    
    	msgPipe := gente.JsonCallbackPipeline{
    		log:    log,
    		router: routerBuilder,
    	}
    
    	srv := websocket.Server{
    		Config:    websocket.Config{},
    		Handler:   wsHandler,
    		Handshake: wsHandshake,
    	}
    
    	http.Handle("/ws", gente.NewConnection(jsonCallbackPipeline{}, log).serveWs)
    
    	log.Fatal(http.ListenAndServe(":8080", handler))
    }

No javascript lib as of yet, but it'll hopefully involve some code generation.
