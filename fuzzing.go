package main

import (
	"fmt"
	"os"
	"bufio"
	"flag"
	"net/http"
	"strings"
	"strconv"
	"time"
)



func search_status_code(status_code_slice []int , status_code_to_compare int , url string){
	for _,c := range status_code_slice {
		if c == status_code_to_compare{
			fmt.Printf("Peticion a la web hecha %s | status_code:%d \n", url,status_code_to_compare)
		}
	}
}

func make_fuzzing(url string , user_agent string , personalized_stc []int){

	if len(personalized_stc) == 0 {
		personalized_stc = []int{200 , 404}
	}
	
	clienteHTTP := &http.Client{}
	
	// Crear una petición GET
	peticion, _ := http.NewRequest("GET", url, nil)
    
    if user_agent != ""{
		peticion.Header.Add("User-Agent" , user_agent)
    }

	// Realizar la petición
	respuesta, err := clienteHTTP.Do(peticion)

	if err != nil {
		fmt.Printf("Porfavor verifique la url %s \n" , url)
	}
	//El defer es para que se ejecute esta opcion asi exista un error en el main
	defer respuesta.Body.Close()

	// Verificar el status_code
	search_status_code(personalized_stc , respuesta.StatusCode , url)
	
}

func main() {
	var wordlist string
	var url string
	var user_agent string
	var stc string
	var personalized_status_code []int
	var delay int
	help := flag.Bool("help", false, "")

	//flags para pedir argumentos
	flag.StringVar(&wordlist , "wordlist" ,"default", "Wordlist para el fuzzing")
	flag.StringVar(&url , "url" , "default", "Url a fuzzear, Ejemplo https://example.com")
	flag.StringVar(&user_agent , "user_agent" , "" ,"User-agent personalizado")
	flag.IntVar(&delay , "delay" , 2 ,"Delay (segundos) entre peticiones")
	flag.StringVar(&stc , "stc" , "" , "Solo mostrara la peticion con el status code que definas")
	flag.Parse()
	
	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if url[len(url) - 1:] == "/"{
		url = url[:len(url) - 1]
		fmt.Println(url)
	}
	
	if url == "" {
		fmt.Println("Porfavor ingrese una url")
		os.Exit(1)
	}

	if stc != "" {
		s := strings.Split(stc , ",")
		for _,f := range s{
			ds , _ := strconv.Atoi(f)
			personalized_status_code = append(personalized_status_code , ds)
		}
	}
	
	fmt.Println(personalized_status_code)

	fmt.Printf("Wordlist: %s \n" , wordlist)

	// Abre el archivo para lectura
	file, err := os.Open(wordlist)

	if err != nil {
		fmt.Printf("Ocurrio un error al abrir el archivo %s \n" , wordlist)
	}

	defer file.Close()

	// Crea un escáner para leer el archivo línea por línea
	scanner := bufio.NewScanner(file)

	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		linea := scanner.Text()
		url_fuzzing := fmt.Sprintf("%s/%s" , url , linea)
		make_fuzzing(url_fuzzing , user_agent , personalized_status_code)
		time.Sleep(time.Duration(delay) * time.Second )
	}
}
