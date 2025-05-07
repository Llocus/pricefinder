function GetUnique(arr: any, index: any) {
    const unique = arr
         .map((e: any) => e[index])
         .map((e: any, i: number, final: any) => final.indexOf(e) === i && i)
        .filter((e: any) => arr[e]).map((e: any) => arr[e]);      
     return unique;
  }
 

  export default GetUnique;