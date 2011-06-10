package main

import (
	"log";
	"os";
	"strconv";
	"strings";
	"bufio";
	"flag";
	"fmt";
	"container/vector";
)

var (
	in *bufio.Reader;
	commands vector.Vector;
	DataStack vector.Vector;
);

var file *string = flag.String("file", "", "Define a Fourth src file to be read");

func main() {
	flag.Parse();
	
	DataStack := new (vector.Vector);
	
	if *file != "" {
		file_to_read, err := os.OpenFile(*file, os.O_RDONLY, 0755); 
		if file_to_read == nil { log.Fatal(err); }
		
		in = bufio.NewReader( file_to_read );
		
		comment := false; //if we're currently in a comment
		
		for {
			dat, err := in.ReadString(' ');
			if err != nil { log.Fatal("can't read file string"); return; }
			
			if dat[0:len(dat)-1] == "(" { 
				comment = true;
			}
			
			if comment == false && dat[0:len(dat)-1] != "" && dat[0:len(dat)-1] != " " {
				parse_forth(dat[0:len(dat)-1], DataStack);
			}
			
			if dat[0:len(dat)-1] == ")"  {
				comment = false;
			}
		}
	} else {
		buf := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("> ");
			read, err := buf.ReadString('\n')
			if err != nil {
				println()
				break
			}
			
			comms := strings.Split(read," ", -1)
			
			for i := 0; i < len(comms); i++ {
				if len(comms) == 0 || comms[i] == "" || comms[i] == " " {
					continue
				}

				parse_forth(strings.TrimSpace(comms[i]), DataStack)
			}
		}
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

func check_stack_size( DataStack *vector.Vector, required int) bool {
	if DataStack.Len() < required {
		log.Fatal("Stack depth is less then " + string(required) + ". Operation is impossible");
		return false;
	} 
	
	return true;
}

func parse_forth(dat string, DataStack *vector.Vector) {
	L := DataStack.Len();
	
	switch strings.TrimSpace(string(dat)) {
		case "":
		case "<cr>":
			return;
		case "t":
			//check the DataStack size using the popped value
			//	if it passes, then the program continues
			minimum := int(DataStack.Pop().(float32));
			if DataStack.Len() < minimum {
				log.Println("DataStack has not enough minerals (values)");
			}
		case ".":
			log.Println(DataStack.Pop());
		case "0SP":
			DataStack.Cut(0, L);
		case ".S":
			log.Println(DataStack);
		case "2/":
			DataStack.Push( DataStack.Pop().(float32) / 2);
		case "2*":
			DataStack.Push( DataStack.Pop().(float32) * 2);
		case "2-":
			DataStack.Push( DataStack.Pop().(float32) - 2);
		case "2+":
			DataStack.Push( DataStack.Pop().(float32) + 2);
		case "1-":
			DataStack.Push( DataStack.Pop().(float32) - 1);
		case "1+":
			DataStack.Push( DataStack.Pop().(float32) + 1);
		case "DUP":
			DataStack.Push( DataStack.Last() );
		case "?DUP":
			if DataStack.Last().(float32) != 0 { DataStack.Push( DataStack.Last().(float32) ); }
		case "PICK":
			number := int(DataStack.Pop().(float32)) ;
			
			if number < L {
				DataStack.Push( DataStack.At(L - 1 - number).(float32) );
			} else {
				log.Fatal("picking out of stack not allowed. Stack Length: " + string(L) + ". Selecting: " + string(number) + ".");
				return;
			}
		case "TUCK":
			DataStack.Insert(L - 2, int(DataStack.Last().(float32)) );
		case "NIP":
			DataStack.Delete(L - 2);
		case "2DROP":
			DataStack.Pop(); DataStack.Pop();
		case "2DUP":
			DataStack.Push(DataStack.At(L - 2));
			DataStack.Push(DataStack.At(DataStack.Len() - 2));
		case "DROP":
			DataStack.Pop();
		case "OVER":
			DataStack.Push(DataStack.At(L - 2));
		case "SWAP":
			l := DataStack.Len();
			DataStack.Swap(l - 2, l - 1);
		case "*":
			num1 := DataStack.Pop().(float32);
			num2 := DataStack.Pop().(float32);
			DataStack.Push( num1 * num2 );				
		case "+":
			num1 := DataStack.Pop().(float32);
			num2 := DataStack.Pop().(float32);
			DataStack.Push( num1 + num2 );
		case "-":
			num1 := DataStack.Pop().(float32);
			num2 := DataStack.Pop().(float32);
			DataStack.Push( num2 - num1 );
		case "/":
			num1 := DataStack.Pop().(float32);
			num2 := DataStack.Pop().(float32);
			DataStack.Push( num2 / num1 );
		case "-ROT":
			DataStack.Swap(L - 1, L - 2);
			DataStack.Swap(L - 2, L - 3);
		case "ROT":
			DataStack.Swap(L - 3, L - 2);
			DataStack.Swap(L - 2, L - 1);
		case "2OVER":
			DataStack.Push(DataStack.At(L - 4));
			DataStack.Push(DataStack.At(DataStack.Len() - 4));
		case "2SWAP":
			DataStack.Swap(L - 4, L - 2);
			DataStack.Swap(L - 3, L - 1);
		case "EMIT":
			log.Println( string([]byte{uint8(DataStack.Last().(float32))}) );
		default:
			val, ok := strconv.Atof32(dat);
			
			if ok == nil {
				DataStack.Push( val );
			} else {
				log.Println(ok);
				log.Fatalln("error, unknown token \""+dat+"\"");
			}
	}
} 

