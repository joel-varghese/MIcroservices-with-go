package main

import "net/http"
import "log"
import "handlers"
import "env"
import "os"
import "time"
import "context"
import "os/signal"

var bindAddress = env.String("BIND_ADDRESS",false,":9090","Bind address for the server")

func main(){
	env.Parse()

	l := log.New(os.Stdout,"product-api",log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/",ph)


	s := &http.Server{
		Addr: *bindAddress,
		Handler: sm,
		ErrorLog: l,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}
	go func(){
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n",err)
			os.Exit(1)
		}
	}()
	c := make(chan os.Signal,1)
	signal.Notify(c,os.Interrupt)
	signal.Notify(c,os.Kill)

	sig := <- c
	log.Println("Got signal", sig)
	

	ctx, _ := context.WithTimeout(context.Background(),30*time.Second)
	s.Shutdown(ctx)
}
