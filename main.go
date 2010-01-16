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
);

var file *string = flag.String("file", "", "Define a Fourth src file to be read");

func main() {
	flag.Parse();
	
	DataStack := vector.New(0);
	
	if *file != "" {
		file_to_read, err := os.Open(*file, os.O_RDONLY, 0755); 
		if file_to_read == nil { log.Stderr(err); }
		
		in = bufio.NewReader( file_to_read );
		
		comment := false; //if we're currently in a comment
		
		for {
			dat, err := in.ReadString(' ');
			if err != nil { log.Stdout("ok"); return; }
			
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
	/**/}

func check_stack_size( DataStack *vector.Vector, required int) bool {
	if DataStack.Len() < required {
		log.Stderr("Stack depth is less then " + string(required) + ". Operation is impossible");
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
			minimum := int(DataStack.Pop().(float));
			if DataStack.Len() < minimum {
				log.Stderr("DataStack has not enough minerals (values)");
			}
		case ".":
			log.Stdout(DataStack.Pop());
		case "0SP":
			DataStack.Cut(0, L);
		case ".S":
			log.Stdout(DataStack);
		case "2/":
			DataStack.Push( DataStack.Pop().(float) / 2);
		case "2*":
			DataStack.Push( DataStack.Pop().(float) * 2);
		case "2-":
			DataStack.Push( DataStack.Pop().(float) - 2);
		case "2+":
			DataStack.Push( DataStack.Pop().(float) + 2);
		case "1-":
			DataStack.Push( DataStack.Pop().(float) - 1);
		case "1+":
			DataStack.Push( DataStack.Pop().(float) + 1);
		case "DUP":
			DataStack.Push( DataStack.Last() );
		case "?DUP":
			if DataStack.Last().(float) != 0 { DataStack.Push( DataStack.Last().(float) ); }
		case "PICK":
			number := int(DataStack.Pop().(float)) ;
			
			if number < L {
				DataStack.Push( DataStack.At(L - 1 - number).(float) );
			} else {
				log.Stderr("picking out of stack not allowed. Stack Length: " + string(L) + ". Selecting: " + string(number) + ".");
				return;
			}
		case "TUCK":
			DataStack.Insert(L - 2, int(DataStack.Last().(float)) );
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
			num1 := DataStack.Pop().(float);
			num2 := DataStack.Pop().(float);
			DataStack.Push( num1 * num2 );				
		case "+":
			num1 := DataStack.Pop().(float);
			num2 := DataStack.Pop().(float);
			DataStack.Push( num1 + num2 );
		case "-":
			num1 := DataStack.Pop().(float);
			num2 := DataStack.Pop().(float);
			DataStack.Push( num2 - num1 );
		case "/":
			num1 := DataStack.Pop().(float);
			num2 := DataStack.Pop().(float);
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
			log.Stdout( string([]byte{uint8(DataStack.Last().(float))}) );
		default:
			val, ok := strconv.Atof(dat);
			
			if ok == nil {
				DataStack.Push( val );
			} else {
				log.Stderr(ok);
				log.Stderr("error, unknown token \""+dat+"\"");
			}
	}
} 

