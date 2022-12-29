package main

import (
	"context"
	"fmt"
	"gin_frame/pkg/setting"
	"gin_frame/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化路由
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpServerPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//certFilePath := filepath.Join("./cert", "public.crt")
	//keyFilePath := filepath.Join("./cert", "private.key")
	// 优雅重启
	go func() {

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//if err := s.ListenAndServeTLS(certFilePath, keyFilePath); err != nil && err != http.ErrServerClosed {
			//log.Fatalf("listen: %s\n", err)
			log.Printf("Listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
