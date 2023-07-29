const throttle = <F extends (...args: Parameters<F>) => ReturnType<F>>( 
    fn: F, 
    delay: number, 
    options: { leading?: boolean; trailing?: boolean } = { 
      leading: true, 
      trailing: true, 
    } 
  ) => { 
    let timer: ReturnType<typeof setTimeout> | undefined 
   
    return (...args: Parameters<F>): void => { 
      if (timer !== undefined) return 
   
      options.leading && fn(...args) 
   
      timer = setTimeout(() => { 
        timer = undefined 
        options.trailing && fn(...args) 
      }, delay) 
    } 
  } 
  
export default throttle