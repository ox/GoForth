package main

import (
	"log";
	"os";
	"strconv";
	"strings";
	"bufio";
	"flag";
	"container/vector";
)

var (
	in *bufio.Reader;
	commands vector.Vector;
	DataStack vector.Vector;
	debug bool;
);

type Object struct {
	name string;
	data string;
}

func NewObject(name string, data string) *Object {
	return &Object{name, data};
}

func (o *Object) String() string {
	return o.data;
}

var file *string = flag.String("file", "", "Define a Fourth src file to be read");

func main() {
	debug := true;
	flag.Parse();
	
	DataStack := vector.New(0);
	
	if *file != "" {
		file_to_read, err := os.Open(*file, os.O_RDONLY, 0755); 
		if file_to_read == nil { log.Stderr(err); }
		
		in = bufio.NewReader( file_to_read );
		
		for {
			dat, err := in.ReadString(' ');
			if err != nil { log.Stderr(err); return; }
			
			parse_forth(dat[0:len(dat)-1], DataStack);
			}
		
	} 
	
	/*
	Experimental, session type operation (think python shell)
	
	
	else {
		in = bufio.NewReader(os.Stdin);
		for {
			dat, err := in.ReadString(' ');
			if err != nil { log.Stderr(err); return; }
			
			parse_forth(dat[0:len(dat)-1], DataStack);
			}
	}
	/**/
	
	if debug { log.Stdout(DataStack); }
	log.Stdout("ok\n");
}

func check_stack_size( DataStack *vector.Vector, required int) bool {
	if DataStack.Len() < required {
		log.Stderr("Stack depth is less then " + string(required) + ". Operation is impossible");
		return false;
	} 
	
	return true;
}

func parse_forth(dat string, DataStack *vector.Vector) {
	switch strings.TrimSpace(string(dat)) {
		case "":
		case "<cr>":
			break;
		case "*":
			if check_stack_size(DataStack, 2) {
				num1, _ := strconv.Atof(DataStack.Pop().(string));
				num2, _ := strconv.Atof(DataStack.Pop().(string));
				DataStack.Push( strconv.Ftoa(num1 * num2, 'f', -1)  );
			} else {
				log.Stderr("error at " + string(dat));
				log.Stderr(DataStack);
				break;
			}
			
		case ".":
			if check_stack_size(DataStack, 1) {
				log.Stdout(DataStack.Pop());
			} else {
				log.Stderr("error at " + string(dat));
				log.Stderr(DataStack);
				break;
			}
			
		case "+":
			if check_stack_size(DataStack, 2) {
				num1, _ := strconv.Atof(DataStack.Pop().(string));
				num2, _ := strconv.Atof(DataStack.Pop().(string));
				DataStack.Push( strconv.Ftoa(num1 + num2, 'f', -1) );
			} else {
				log.Stderr("error at " + string(dat));
				log.Stderr(DataStack);
				break;
			}
			
		case "-":
			if check_stack_size(DataStack, 2) {
				num1, _ := strconv.Atof(DataStack.Pop().(string));
				num2, _ := strconv.Atof(DataStack.Pop().(string));
				DataStack.Push( strconv.Ftoa(num1 - num2, 'f', -1) );
			} else {
				log.Stderr("error at " + string(dat));
				log.Stderr(DataStack);
				break;
			}
			
		case "/":
			if check_stack_size(DataStack, 2) {
				num1, _ := strconv.Atof(DataStack.Pop().(string));
				num2, _ := strconv.Atof(DataStack.Pop().(string));
				DataStack.Push( strconv.Ftoa(num1 / num2, 'f', -1) );
			} else {
				log.Stderr("error at " + string(dat));
				log.Stderr(DataStack);
				break;
			}
		default:
			_, ok := strconv.Atof(dat);
			
			if ok == nil {
				DataStack.Push( dat );
			} else {
				log.Stderr(ok);
				log.Stderr ("error, unknown token \""+string(dat)+"\"");
				}
	}
	
	if debug { log.Stdout( DataStack ); }
}