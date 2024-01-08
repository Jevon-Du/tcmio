library(alakazam)
library(shazam)
library(dplyr)
library(ggplot2)

# the path of clone typed file with assigned germline

for (i in c('5_2',20)) {
    t_file = paste('/data/data/clones/CM', i, '_train_clone-pass_germ-pass_merged.tsv', sep='')
    db_t = read.csv(t_file, header = TRUE, sep = '\t')
    db_obs_t <- observedMutations(db_t,
                                sequenceColumn='sequence_alignment',
                                germlineColumn='germline_alignment_d_mask',
                                regionDefinition=NULL,
                                frequency=TRUE,
                                nproc=5)
    t_ofile = paste('/data/data/clones/CM', i, '_train_clone-pass_germ-pass_mutation.tsv', sep='')
    write.table(db_obs_t, t_ofile, sep='\t', quote=FALSE, col.names=TRUE, row.names=FALSE)
    print(paste('Done: ', t_ofile, sep=''))
    
    if (i!=0){
        v_file = paste('/data/data/clones/CM', i, '_valid_clone-pass_germ-pass_merged.tsv', sep='')
        db_v = read.csv(v_file, header = TRUE, sep = '\t')
        db_obs_v <- observedMutations(db_v, 
                                sequenceColumn='sequence_alignment',
                                germlineColumn='germline_alignment_d_mask',
                                regionDefinition=NULL,
                                frequency=TRUE, 
                                nproc=5)
        v_ofile = paste('/data/data/clones/CM', i, '_valid_clone-pass_germ-pass_mutation.tsv', sep='')
        write.table(db_obs_v, v_ofile, sep='\t', quote=FALSE, col.names=TRUE, row.names=FALSE)
        print(paste('Done: ', v_ofile, sep=''))
    }
}
