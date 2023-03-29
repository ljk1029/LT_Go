package main 
  
import "fmt"
  
func add(a1, a2 int) int { 
    res := a1 + a2 
    fmt.Println("Result: ", res) 
    return 0 
}   
  //匿名函数作为参数传递  
 func GFG(i func(p, q string)string){ 
    fmt.Println(i ("Geeks", "for"))  
 } 
    
func main() { 
	fmt.Println("Start") 
  
    //多个延迟语句
    //以LIFO顺序执行
    defer fmt.Println("End") 
    defer add(34, 56) 
    defer add(10, 10) 
    value:= func(p, q string) string{ 
    	return p + q + "Geek"
    } 
	GFG(value) 
}