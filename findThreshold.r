library(alakazam)
library(ape)
library(shazam)
# library(pryr)
# library(ggplot2)

find_opt_thres <- function(db_tsv){
    db = read.csv(db_tsv, header = TRUE, sep='\t')
    dist_ham <- distToNearest(db, sequenceColumn="junction", vCallColumn="v_call", jCallColumn="j_call",
                            model="ham", symmetry="avg", normalize="len", nproc=8)
    m="density"
    # Find threshold using the "gmm" method with optimal threshold
    dist_ham_nona <- subset(dist_ham, !is.na(dist_nearest))$dist_nearest
    
    output <- findThreshold(dist_ham_nona)
    
    # Find threshold using gmm method
    if (is.na(output@threshold)) {
        m="gmm"
        output <- findThreshold(dist_ham_nona, method="gmm", model="gamma-gamma")
    }
    
    if (is.null(output)) {
        res = c(db_tsv, NA, m)
    } else {
        res = c(db_tsv, round(output@threshold, digits=3), m)
    }
    print(res)
    return(res)
}
  # p1 <- ggplot(subset(dist_ham, !is.na(dist_nearest)),
  #              aes(x=dist_nearest)) + theme_bw() +
  #              xlab("Hamming distance") + ylab("Count") +
  #              scale_x_continuous(breaks=seq(0, 1, 0.1)) +
  #              geom_histogram(color="white", binwidth=0.02) +
  #              geom_vline(xintercept=output@threshold, color="firebrick", linetype=2)
  # plot(p1)
    
header = c('file_in_docker', 'thresold', 'method')
# create an empty data.frame
res_df = as.data.frame(matrix(nrow=0,ncol=3))
colnames(res_df) <- header

# CM00_file = '/data/data/cell_pass/CM00_train_yao_fmt7.tsv'
# res_t = find_opt_thres(CM00_file)
# res_df[nrow(res_df)+1,] = res_t
# write.table(res_df, '/data/data/threshold_CM.csv', sep=',', append=TRUE, quote=FALSE, col.names=TRUE, row.names=FALSE)

for (i in c('7_new')) {

    t_file = paste('/data/data/cell_pass/CM', i, '_train.tsv', sep='')
    res_t = find_opt_thres(t_file)
    res_df[nrow(res_df)+1,] = res_t
    
    if (i!=0){
        v_file = paste('/data/data/cell_pass/CM', i, '_valid.tsv', sep='')
        if (file.exists(v_file)){
            res_v = find_opt_thres(v_file)
            res_df[nrow(res_df)+1,] = res_v
        } else {
            tp_file = paste('/data/data/cell_pass/CM', i, '_valid_tp.tsv', sep='')
            tp_res_v = find_opt_thres(tp_file)
            res_df[nrow(res_df)+1,] = tp_res_v
            
            tn_file = paste('/data/data/cell_pass/CM', i, '_valid_tn.tsv', sep='')
            tn_res_v = find_opt_thres(tn_file)
            res_df[nrow(res_df)+1,] = tn_res_v
        }            
    }
}

write.table(res_df, '/data/data/threshold_CM_20220809.csv', sep=',', append=TRUE, quote=FALSE, col.names=FALSE, row.names=FALSE)

# for train+valid data
# for (i in c('5_2',20)) {
#     ifile = paste('/data/data/cell_pass/CM', i, '_fixed.tsv', sep='')
#     if (!file.exists(ifile)){
#         ifile = paste('/data/data/cell_pass/CM', i, '.fmt7', sep='')
#     }
#     res_t = find_opt_thres(ifile)
#     res_df[nrow(res_df)+1,] = res_t
# }

# write.table(res_df, '/data/data/threshold_CM_20220823.csv', sep=',', quote=FALSE, col.names=TRUE, row.names=FALSE)