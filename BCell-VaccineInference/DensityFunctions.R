## Load required packages ##
library(VGAM)
library(ismev)
library(MASS)

#################################################


source('ShortFunctions.R')

#################Discrete GPD density###############



dgpd_discrete<-function(x,location,scale,shape,log){ #discrete version
    
    if(length(x)==1){ if( is.na(x)){return(NA)} }
    
    if(log){	
	return(log(pgpd(x+1,location,scale,shape) - pgpd(x,location,scale,shape)))
    }else{
	return(pgpd(x+1,location,scale,shape)-pgpd(x,location,scale,shape))
    }

}

################################################


############Mixture model density###################



Density_NB_GPD<-function(x,nb,gpd,weights){ # Density of negative binomial and gpd mixture model; used on data with 0s extracted for estimating threshold parameter; weights = single value on NEGATIVE BINOMIAL
	
    dens.gpd = dgpd_discrete(x,location=gpd$threshold,scale=gpd$mle[1],shape=gpd$mle[2],log=T)
    dens.nb = dnbinom(x,size=nb$estimate[1],mu=nb$estimate[2],log=T)
    return(sum(mapply(log_mm,dens.nb,dens.gpd,weights)))

}





Density<-function(x,nb_1,nb_2,gpd_thresh,gpd_1,gpd_2,weights,component=0){ #density of delta 0/negative binomial/gpd mixture model; input parameters specified separately - should be of a form compatible with the data, x, e.g. vectors of multiples of the same length; component provides option to return individual likelihoods for different components
  
	dens.1 = log(as.numeric(x==0))
	dens.gpd = dgpd_discrete(x,location=gpd_thresh,scale=gpd_1,shape=gpd_2,log=T)
	dens.nb = dnbinom(x,size=nb_1,mu=nb_2,log=T)
    
	if(is.matrix(x)){
		dens.1=matrix(dens.1,nrow=nrow(x))
    	}else if(is.array(x)){
		dens.1=array(dens.1,dim=dim(x))
		dens.gpd=array(t(apply(x,1,function(y1)dgpd_discrete(y1,location=gpd_thresh,scale=gpd_1,shape=gpd_2,log=T))),dim=dim(x))
		dens.nb=array(t(apply(x,1,function(y)dnbinom(y,size=nb_1,mu=nb_2,log=T))),dim=dim(x))
    	}
    
	if(component==0){
		return(sum(log_mm3(dens.1,dens.nb,dens.gpd,weights)))
    	}else if(component==1){
		return(dens.1)
    	}else if(component==2){
		return(dens.nb)
    	}else if(component==3){
		return(dens.gpd)
	}else if(component==4){
		return(log_mm3(dens.1,dens.nb,dens.gpd,weights))
	}else if(component==5){
        	return(list(log(weights[1])+dens.1,log(weights[2])+dens.nb,log(weights[3])+dens.gpd,log_mm3(dens.1,dens.nb,dens.gpd,weights)))
    	}

}	


Observed_Density_classified<-function(x,nb_1,nb_2,gpd_thresh,gpd_1,gpd_2,weights,p_pres,component=0,FirstOnly=FALSE){#Density for model which has seperate components for sequences which are present/not present in a sample; requires x to be a time x subject matrix and input model parameters to be of the same dimensions
#FirstOnly if only returning density of first time point - still need to consider all data to determine absence from sample/time point
#conditional on a classification - used to optimise parameters
	
	N_sub = ncol(x)
	Pres = which(colSums(x)!=0)
	N_pres=length(Pres)
	
	Density_pres = Density(x[,Pres],nb_1[,Pres],nb_2[,Pres],gpd_thresh[,Pres],gpd_1[,Pres],gpd_2[,Pres],weights,component) 
	
	
	if(FirstOnly){
		if(component==5){
			if(length(Pres)>1){
				Density_pres=lapply(Density_pres,function(x)x[1,])
			}else{
				Density_pres=lapply(Density_pres,function(x)x[1])
			}
		}else{
			if(length(Pres)>1){
				Density_pres=Density_pres[1,]
			}else{
				Density_pres=Density_pres[1]
			}
		}

	}

        
	if(component==0){
            Density_pres = Density_pres + N_pres*log(p_pres)
	    Density_notPres = (N_sub - N_pres)*log(1-p_pres)

	    return(Density_pres + Density_notPres)
        }else{
            return(Density_pres)
        }
}




Observed_Density_allclasses<-function(x,nb_1,nb_2,gpd_thresh,gpd_1,gpd_2,weights_bg,weights_ind,p_pres_bg,p_pres_ind){#Density for model which has seperate components for sequences which are present/not present in a sample; requires x to be a time x subject matrix and input model parameters to be of the same dimensions
#Not conditional on a classification, all classifications considered with some weight - used to optimise classifications/allocations/assignments
	N_sub = ncol(x)
	Pres = which(colSums(x)!=0)
	N_pres=length(Pres)
	
	
	Density_pres_BG = Density(x[,Pres,drop=F],nb_1[,Pres,drop=F],nb_2[,Pres,drop=F],gpd_thresh[,Pres,drop=F],gpd_1[,Pres,drop=F],gpd_2[,Pres,drop=F],weights_bg,component=4) 
        Density_pres_IND = Density(x[,Pres,drop=F],nb_1[,Pres,drop=F],nb_2[,Pres,drop=F],gpd_thresh[,Pres,drop=F],gpd_1[,Pres,drop=F],gpd_2[,Pres,drop=F],weights_ind,component=4) 
        
        Density_pres_bg = sum(Density_pres_BG) + N_pres*log(p_pres_bg)
	Density_notPres_bg = (N_sub - N_pres)*log(1-p_pres_bg)

        Density_pres_bg_ind = sum(Density_pres_IND) + N_pres*log(p_pres_bg)
        Density_notPres_bg_ind = (N_sub - N_pres)*log(1-p_pres_bg)

        Density_pres_ind = sum(Density_pres_BG[1,])+sum(Density_pres_IND[-1,]) + N_pres*log(p_pres_ind)        
        Density_notPres_ind = (N_sub - N_pres)*log(1-p_pres_ind)
    

        return(c(Density_pres_bg + Density_notPres_bg, Density_pres_ind + Density_notPres_ind, Density_pres_bg_ind+Density_notPres_bg_ind))
}



########################################################



