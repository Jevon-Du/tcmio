## Load required packages ##
library(VGAM)
library(ismev)
library(MASS)



source('ShortFunctions.R')
source('DensityFunctions.R')



#############Threshold estimation#####################



thresh.est<-function(thresh,data){	#function which can be maximised over threshold parameter; Need to pre-filter 0s out of data.
     
 	if(sum(as.vector(data)>thresh)<=1){
   		return(Inf)
 	}else{	
        	GPD<-gpd.fit(as.vector(data),threshold=thresh,method="L-BFGS-B")
		NB<-fitdistr(data[data<=thresh],"negative binomial",start=list(mu=mean(data[data<=thresh]),size=1),lower=c(.000001,0.000001))
		res = Density_NB_GPD(data,NB,GPD,min(1-(sum(data>thresh)/length(data)),.999))
        	return(-res)	
	}
}




###########Parameter estimation only (no classification)##############

Fit_initial_params<-function(Data,p_bg=.85,p_ind=.1,w_bg=c(.59,.309,0.0001),w_ind=c(.1,.2,.7),allocations=NULL){		#Data should be a cluster x time x subject array
#Fits independent set of parameters to data from each sample and each time point.
   
	Nclust=dim(Data)[1]
	Ntime=dim(Data)[2]
	Nsamp=dim(Data)[3]

   	thresholds=matrix(nrow=Ntime,ncol=Nsamp)
    	gpd_1=matrix(nrow=Ntime,ncol=Nsamp)
    	gpd_2=matrix(nrow=Ntime,ncol=Nsamp)
    	nb_1=matrix(nrow=Ntime,ncol=Nsamp)
    	nb_2=matrix(nrow=Ntime,ncol=Nsamp)


	joint_weights=p_bg*w_bg + p_ind+w_ind

    
    
    	for(j in 1:Nsamp){
		
		Data1<-Data[,,j]
        	Data1<-Data1[rowSums(Data1)!=0,] #remove clusters absent from sample j
        
		for(i in 1:Ntime){
            		
			print(paste("Time point: ",i,sep=""))
            		print(paste("Sample: ",j,sep=""))
            	
			data=Data1[,i]
			non_zero_data<-c(data[data!=0],rep(0,ceiling(joint_weights[2]*sum(data==0)))) 	#filter 0s accounted for by point mass

			if(max(non_zero_data)==1){

				print('WARNING: maximum cluster size of 1')
    				thresh=2
    				thresholds[i,j]=2
    				GPD<-list(threshold=2,mle=c(1,1))

			}else{
	        
			    	range=2:100              #Brute force alternative to optimize; TO DO consider range values....
            			ll_thresh=c()
	    			
				for(m in range){
					temp=try(thresh.est(m,non_zero_data))
					if(is.numeric(temp[1])){ll_thresh=c(ll_thresh,temp[1])}else{ll_thresh=c(ll_thresh,Inf)}
	    			}
            
				thresholds[i,j]=range[which.min(ll_thresh)]
            			thresh=range[which.min(ll_thresh)]
           
            			GPD<-gpd.fit(data,threshold=thresh)
         	
			}

			if(is.null(allocations)){
				NB_data<-c(non_zero_data[non_zero_data<=thresh],sample(non_zero_data[non_zero_data>thresh],floor(sum(non_zero_data>thresh)*joint_weights[2])))
			}else{

				NB_data<-Data[allocations==1,i,j]
				NB_data<-NB_data[NB_data>0]
				NB_data<-c(NB_data,rep(0,ceiling(sum(NB_data==0)*joint_weights[2])))	#filter 0s accounted for by point mass
			}
 		
			NB<-fitdistr(NB_data,"negative binomial",start=list(mu=mean(NB_data),size=1),lower=0.000001)
            
			gpd_1[i,j]=GPD$mle[1]
            		gpd_2[i,j]=GPD$mle[2]
            		nb_1[i,j]=NB$estimate[1]
            		nb_2[i,j]=NB$estimate[2]

    		}
	}

	ll<-apply(Data,1,Observed_Density_allclasses,nb_1,nb_2,thresholds,gpd_1,gpd_2,w_bg,w_ind,p_pres_bg=1/dim(Data)[3],p_pres_ind=2/dim(Data)[3])
print(nb_2)
	
	
	return(list(nb_1,nb_2,thresholds,gpd_1,gpd_2,sum(ll[1,])))

}


###################################################################
