## Load required packages ##
library(VGAM)
library(ismev)
library(MASS)



################ Log-likelihood for mixture model###################

log_mm<-function(x,y,weights){   # log(w*exp(f) + (1-w)*exp(g))

    	x1<-which(x>y)
    	res<-x
   	if(length(weights)==1){weights=c(weights,1-weights)}
    
	if(length(x1)>0){
		
		res[x1]<- x[x1] + log(weights[1] + weights[2]*exp(y[x1]-x[x1]))
		if(length(x1)<length(x)){ res[-x1]<- y[-x1] + log(weights[1]*exp(x[-x1]-y[-x1]) + weights[2]) }
    	
	}else{
		res<- y + log(weights[1]*exp(x-y) + weights[2])
    	}
    	
	return(res)

}




log_mm3<-function(x,y,z,weights){ # log(w_1*exp(f) + w_2*exp(g) + w_3*exp(h))
	
    res1<-log_mm(x,y,weights[1:2])
    res2<-log_mm(res1,z,c(1,weights[3]))
    return(res2)
}

######################################################################


