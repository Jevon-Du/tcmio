## Load required packages ##
library(VGAM)
library(ismev)
library(MASS)


###############################################

source('ShortFunctions.R')
source('DensityFunctions.R')
source('ParameterEstimation.R')


#################################################




#####EM Algorithm for model fitting #########################################


EM<-function(Data, max_iter,p_bg=.85,p_ind=.1,w_bg=c(.99,.009,0.001),p_pres_bg=.25,p_pres_ind=.5,w_ind=c(.1,.2,.7),Params=NULL){		#Data should be a cluster x time x subject array
#Fits independent set of parameters to data from each sample and each time point.
#Fits model for multiple samples; cluster either present in sample or not present in sample

	ll_track<-vector(length=max_iter)
    
    	Nclust=dim(Data)[1]
    	Ntime=dim(Data)[2]
    	Nsamp=dim(Data)[3]
    
	if(is.null(Params)){Params<-Fit_initial_params(Data=Data,p_bg=p_bg,p_ind=p_ind,w_bg=w_bg,w_ind=w_ind)}   
    
	nb_1=Params[[1]]
    	nb_2=Params[[2]]
   	thresholds=Params[[3]]
   	gpd_1=Params[[4]]
    	gpd_2=Params[[5]]

    	assignments<-rep(0,Nclust)
    	count=0	

    	ll=apply(Data,1,Observed_Density_allclasses,nb_1=nb_1,nb_2=nb_2,gpd_thresh=thresholds,gpd_1=gpd_1,gpd_2=gpd_2,weights_bg=w_bg,weights_ind=w_ind,p_pres_bg=p_pres_bg,p_pres_ind=p_pres_ind)
    
	LL<-rbind(log(p_bg)+ll[1,],log(p_ind)+ll[2,],log(1-p_bg-p_ind)+ll[3,])

    	assignments_new=apply(rbind(log(p_bg)+ll[1,],log(p_ind)+ll[2,],log(1-p_bg-p_ind)+ll[3,]),2,which.max)
    
	
	while(count<max_iter && sum(assignments_new==assignments)!=length(assignments)){
      

   		count = count + 1
		assignments<-assignments_new
   
		if(sum(assignments==3)>0){ 
        	
			W_IND1=apply(Data[assignments==2,-1,,drop=F],1,Observed_Density_classified,nb_1[-1,],nb_2[-1,],thresholds[-1,],gpd_1[-1,],gpd_2[-1,],w_ind,component=5,p_pres=p_pres_ind)
        		W_IND2=apply(Data[assignments==3,,,drop=F],1,Observed_Density_classified,nb_1=nb_1,nb_2=nb_2,gpd_thresh=thresholds,gpd_1=gpd_1,gpd_2=gpd_2,w=w_ind,component=5,p_pres=p_pres_ind)
        		ind_baseline=unlist(c(lapply(W_IND1,function(x)x[[4]]),lapply(W_IND2,function(x)x[[4]]))) #all components (individual for each observation)
        		ind_1=unlist(c(lapply(W_IND1,function(x)x[[1]]),lapply(W_IND2,function(x)x[[1]])))       #1st component
        		ind_2=unlist(c(lapply(W_IND1,function(x)x[[2]]),lapply(W_IND2,function(x)x[[2]])))      #2nd
        		ind_3=unlist(c(lapply(W_IND1,function(x)x[[3]]),lapply(W_IND2,function(x)x[[3]])))        #3rd
		
		}else{
		
			W_IND1=apply(Data[assignments==2,-1,,drop=F],1,Observed_Density_classified,nb_1[-1,],nb_2[-1,],thresholds[-1,],gpd_1[-1,],gpd_2[-1,],w_ind,component=5,p_pres=p_pres_ind)
			ind_baseline=unlist(lapply(W_IND1,function(x)x[[4]])) #all components (individual for each observation)
        		ind_1=unlist(lapply(W_IND1,function(x)x[[1]]))       #1st component
        		ind_2=unlist(lapply(W_IND1,function(x)x[[2]]))      #2nd
        		ind_3=unlist(lapply(W_IND1,function(x)x[[3]]))        #3rd
			
		}	

		w_ind[1] = sum(exp(ind_1-ind_baseline))/length(ind_baseline)
		w_ind[2] = sum(exp(ind_2-ind_baseline))/length(ind_baseline)
		w_ind[3] = sum(exp(ind_3-ind_baseline))/length(ind_baseline)
	
        	
		
		W_BG1=apply(Data[assignments==2,,],1,Observed_Density_classified,nb_1,nb_2,thresholds,gpd_1,gpd_2,w_bg,component=5,p_pres=p_pres_bg,FirstOnly=T)
        	W_BG2=apply(Data[assignments==1,,],1,Observed_Density_classified,nb_1,nb_2,thresholds,gpd_1,gpd_2,w_bg,component=5,p_pres=p_pres_bg)
        
        	bg_baseline= unlist(c(lapply(W_BG1,function(x)x[[4]]),lapply(W_BG2,function(x)x[[4]])))        
        	bg_1= unlist(c(lapply(W_BG1,function(x)x[[1]]),lapply(W_BG2,function(x)x[[1]])))        
        	bg_2= unlist(c(lapply(W_BG1,function(x)x[[2]]),lapply(W_BG2,function(x)x[[2]])))        
        	bg_3= unlist(c(lapply(W_BG1,function(x)x[[3]]),lapply(W_BG2,function(x)x[[3]])))        

        	w_bg[1] = sum(exp(bg_1-bg_baseline))/length(bg_baseline)
		w_bg[2] = sum(exp(bg_2-bg_baseline))/length(bg_baseline)
		w_bg[3] = sum(exp(bg_3-bg_baseline))/length(bg_baseline)




       		p_pres_bg=sum(apply(Data[assignments==1|assignments==3,,],1,function(x)sum(colSums(x)!=0)))/(sum(assignments==1|assignments==3)*Nsamp)
       		p_pres_ind=sum(apply(Data[assignments==2,,],1,function(x)sum(colSums(x)!=0)))/(sum(assignments==2)*Nsamp)

		p_bg<-max(sum(assignments==1)/length(assignments),.01)
		p_ind<-max(sum(assignments==2)/length(assignments),.01)
   		if(p_bg+p_ind>=1){p_bg=p_bg-.01}  


		Params<-Fit_initial_params(Data=Data,p_bg=p_bg,p_ind=p_ind,w_bg=w_bg,w_ind=w_ind,allocations=assignments) 


		nb_1=Params[[1]]
		nb_2=Params[[2]]
		thresholds=Params[[3]]
		gpd_1=Params[[4]]
		gpd_2=Params[[5]]


        	ll=apply(Data,1,Observed_Density_allclasses,nb_1=nb_1,nb_2=nb_2,gpd_thresh=thresholds,gpd_1=gpd_1,gpd_2=gpd_2,weights_bg=w_bg,weights_ind=w_ind,p_pres_bg=p_pres_bg,p_pres_ind=p_pres_ind)


   		LL<-rbind(log(p_bg)+ll[1,],log(p_ind)+ll[2,],log(1-p_bg-p_ind)+ll[3,])
   		LL<-rbind(log(p_bg)+ll[1,],log(p_ind)+ll[2,],log(1-p_bg-p_ind)+ll[3,])


      		assignments_new=apply(rbind(log(p_bg)+ll[1,],log(p_ind)+ll[2,],log(1-p_bg-p_ind)+ll[3,]),2,which.max)

		ll_track[count]=sum(LL[cbind(assignments_new,1:length(assignments_new))])
	
        	print(count)
		print(p_bg)
		print(p_ind)
		print(w_bg)
		print(w_ind)
		print(p_pres_bg)
        print(p_pres_ind)
        print(sum(assignments_new==assignments))
		print(sum(assignments_new==2))
    
    
    }
    return(list(assignments_new,w_bg,p_bg,p_pres_bg,w_ind,p_ind,p_pres_ind,ll,ll_track,count))
}


##################################################################################################################


