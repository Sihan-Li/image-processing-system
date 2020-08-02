package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"os"
	"proj2/png"
	"sync"
)

type Command struct {
	InPath   	string 	`json:"inPath"`
	OutPath     string	`json:"outPath"`
	Effects		[]string `json:"effects"`
}

func writer(img *png.Image, outPath string){
	img.Save(outPath)
}

func sequential(){
	dec := json.NewDecoder(os.Stdin)
	for {
		var command Command
		if err := dec.Decode(&command); err != nil {
			log.Println(err)
			return
		}
		img, _ := png.Load(command.InPath)
		img.InitTmp()
		numberOfOperations := len(command.Effects)
		outPath := command.OutPath
		for i := 0; i < numberOfOperations; i++ {
			if command.Effects[i] == "S" {
				img.ConvertImage(img.Sharpen())
			} else if command.Effects[i] == "E" {
				img.ConvertImage(img.Edge_detection())
			} else if command.Effects[i] == "B" {
				img.ConvertImage(img.Blur())
			} else {
				img.Grayscale()
			}
		}
		writer(img,outPath)
	}
}

func reader(done chan string,listOfTasks *[]Command,m *sync.Mutex){
	dec := json.NewDecoder(os.Stdin)
	for{
		var command Command
		m.Lock()
		err := dec.Decode(&command)
		if  err != nil {
			log.Println(err)
			m.Unlock()
			break
		}
		*listOfTasks = append(*listOfTasks,command)
		m.Unlock()
	}
	done <- "read task finished"
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return quotient,remainder
}
func handleParallelTask(img *png.Image,command Command,numberOfThreads int,lengthOfNormalPart int){
	for _,operation := range command.Effects{
		if operation == "S" {
			img.DivideImage(img.Sharpen(),numberOfThreads,lengthOfNormalPart)
		} else if operation == "E" {
			img.DivideImage(img.Edge_detection(),numberOfThreads,lengthOfNormalPart)
		} else if operation == "B" {
			img.DivideImage(img.Blur(),numberOfThreads,lengthOfNormalPart)
		} else if operation == "G" {
			img.DivideGrayscale(numberOfThreads,lengthOfNormalPart)
		}
	}
}
func pipeLine(done chan string,numberOfThreads int,command Command){
	img, _ := png.Load(command.InPath)
	img.InitTmp()
	outPath := command.OutPath
	bounds := img.Out.Bounds()
	width := bounds.Max.X
	lengthOfNormalPart, _ := divmod(width,numberOfThreads)

	handleParallelTask(img,command,numberOfThreads,lengthOfNormalPart)
	writer(img,outPath)
	done <- "One image is finished"
}


func main() {
	if len(os.Args) == 1{
		sequential()
	}else{
		pPtr := flag.Int("p", 0, "a string")
		flag.Parse()
		numOfThreads := *pPtr
		number_Of_Readers := int(math.Ceil(float64(numOfThreads) * (1.0 / 5.0)))

		//We input the chan Command into the reader and load a command onto it
		//readCommands := make(chan Command,number_Of_Readers)
		readDone := make(chan string,number_Of_Readers )
		var listOfTasks []Command
		var mutex sync.Mutex
		for i:=0; i < number_Of_Readers; i++{
			go reader(readDone,&listOfTasks,&mutex)
		}
		for i:=0 ; i < number_Of_Readers ; i++{
			<-readDone
		}


		pipeLineDone := make(chan string, len(listOfTasks))
		//We have a list of commands and then do operations
		for _,c := range listOfTasks{
			go pipeLine(pipeLineDone,numOfThreads,c)
		}
		for i:=0; i < len(listOfTasks); i++{
			<- pipeLineDone
		}
	}
}
