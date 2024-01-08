# %%
import csv
from datetime import datetime
from pathlib import Path
from typing import List, Tuple, Dict, Union, Optional

import numpy as np
import pandas as pd
from prettytable import PrettyTable

from utils import *


###################### 配置 #######################
# path in ubuntu
ROOT = Path('/data/personal/jiewen.du/CM_new')
SCRIPT = ROOT / 'script'
DATA_HOME = ROOT / 'data'
DATA_DIR = DATA_HOME / 'ymc_copy' # 姚盟成处理后的数据文件
CELL_PASS_DIR = DATA_HOME / 'cell_pass' # 过滤掉非B cells, CDR3为空，非productive BCR细胞的数据文件
FMT19_DIR = DATA_HOME / 'fmt19' # igblast对BCR序列按照AIRR格式进行注释后的结果文件
THRESHOLD_FILE = DATA_HOME/'threshold_CM_20220809.csv'

# path in docker container
CONTAINER = 'immcantation_cm'
WORK_DIR = Path('/data')

CURRENT_DATASET = ['4_new', 43]

VALID_MICE.update({
    'CM4_new': [['CM004-2-2-1']],
    'CM43': [['CM043-2-1']]
})

statics_header = ['dataset', 'All cells', 'Non-B cells', 'No-Germlines', 'B cells', 
                  'B cells train', 'B cells train clone_pass', 'unique NT train clone_pass', 
                  'B cells train clone_pass_germ_pass', 'unique NT train clone_pass_germ_pass', 
                  'B cells valid', 'B cells valid clone_pass', 'unique NT valid clone_pass', 
                  'B cells valid clone_pass_germ_pass', 'unique NT valid clone_pass_germ_pass'
                 ]

STATISTIC = {}

#############################################


print(datetime.today().strftime("%Y-%m-%d %H:%M:%S"))
today = datetime.today().strftime("%Y%m%d%H%M%S")

"""
## Create docker container

```bash
docker run -ti --runtime=runc --name immcantation_cm --workdir /data -v /data/personal/jiewen.du/CM_new:/data immcantation/suite:4.3.0 /bin/bash
```
"""


"""
## Step1: Cell-pass

过滤掉非B cells, CDR3为空，非productive BCR的细胞数据

备注：
1. CM0~CM8: CM*bcr.samples.IGHV.igblast.tsv均为fmt19格式
2. CM9, CM10: CM*bcr.samples.IGHV.igblast.tsv均为fmt7格式  
"""

for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    df0 = pd.read_csv(DATA_DIR/f'{dataset}.TSM_VS_VMM_all.bcr.samples.IGHV.igblast.tsv', sep='\t')

    df = df0.loc[df0['celltype'].isin(BTYPE.keys()), :].copy(True)
    print('Source: {:<6} | Cells: {:<6} | Non-Bcells: {:<6} '.format(dataset, df0.shape[0], df0.shape[0]-df.shape[0]), end='')
    
    df.reset_index(drop=True, inplace=True)
    df['mice_index'] = df['sequence_id'].apply(lambda x: x.split('.')[-1])

    droped_idx1 = df.loc[(df['cdr3']=='')|(df['cdr3']=='None')|(df['cdr3']=='NA')|(df['cdr3']==None)|(df['cdr3'].isna()), :].index
    df.drop(index=droped_idx1, inplace=True)
    df.reset_index(drop=True, inplace=True)

    droped_idx2 = df.loc[df['productive']==False, :].index
    df.drop(index=droped_idx2, inplace=True)

    # 2: fix wrong dtype
    df = fix_wrong_dtype(df)
    df = df.where(pd.notna(df), None)

    df.to_csv(CELL_PASS_DIR/f'{dataset}.fmt7', sep='\t', index=None, na_rep='')

    with open(CELL_PASS_DIR/f'{dataset}.fasta', 'w') as f:
        for row in df.loc[:, ['sequence_id', 'sequence']].itertuples(False):
            f.write(f'>{row[0]}\n')
            f.write(f'{row[1]}\n')
    
    print('| Null CDR3: {:<6} | Non-productives: {:<6} | Cell-pass: {:<6}'.format(len(droped_idx1), len(droped_idx1), df.shape[0]))
    statistic[dataset] = [dataset, df0.shape[0], df0.shape[0]-df.shape[0]]


# %%
"""
## Step2: Run igblastn(v1.17.0)

使用igblast对BCR序列按照AIRR格式进行注释，得到CDR和framework各区域的序列

"""
for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    fasta = CELL_PASS_DIR/f'{dataset}.fasta'
    run_igblast(fasta, fmt=19, threds=40)

"""
## Step3: Check db

过滤igblastn结果中的错误比对序列（germline_alignment为空），输出为*fixed.tsv文件；如果没有比对错误，则不输出*fixed.tsv文件

"""

for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    ifile = CELL_PASS_DIR/f'{dataset}.fmt7'
    df = pd.read_csv(ifile, sep='\t')
    df = fix_wrong_dtype(df)
    
    droped_idx1 = df.loc[df['germline_alignment'].isna(),:].index
    droped_idx2 = df.loc[df['d_germline_start']>=df['d_germline_end'],:].index
    print('Source: {:<10} | Null germline_alignment: {:<6} | d_start > d_end: {:<6} | Passed-cells: {:<6}'.format(dataset, len(droped_idx1), len(droped_idx2), df.shape[0]-len(droped_idx1)))
    statistic[dataset].extend([len(droped_idx1), df.shape[0]-len(droped_idx1)])
    
    if len(droped_idx1)!=0 or len(droped_idx2)!=0:
        df.drop(index=droped_idx1,inplace=True)
        df.to_csv(CELL_PASS_DIR/f'{dataset}_fixed.tsv', sep='\t', index=None, na_rep='')

print(statistic)

# %%
"""
## Step4: split sequences to train and valid

将一个课题的B细胞分为一次流式和二次流式，供后续分别做clonotype聚类和germline推断
"""

for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    ifile = CELL_PASS_DIR/f'{dataset}_fixed.tsv'
    if not ifile.exists():
        ifile = CELL_PASS_DIR/f'{dataset}.fmt7'
    
    df = pd.read_csv(ifile, sep='\t')
    # df['mice_index'] = df['sequence_id'].apply(lambda x: x.split('.')[1])
    df = fix_wrong_dtype(df)
    
    df_tp = df.loc[df['mice_index'].isin(VALID_MICE[dataset][0]), :]
    if len(VALID_MICE[dataset])>1:
        df_tn = df.loc[df['mice_index'].isin(VALID_MICE[dataset][1]), :]
        df_t = df.loc[(~df.index.isin(df_tp.index))&(~df.index.isin(df_tn.index)), :]
        df_tn.to_csv(CELL_PASS_DIR/f'{dataset}_valid_tn.tsv', sep='\t', index=None, na_rep='')
        df_tp.to_csv(CELL_PASS_DIR/f'{dataset}_valid_tp.tsv', sep='\t', index=None, na_rep='')
    else:
        df_t = df.loc[~df.index.isin(df_tp.index), :]
        df_tp.to_csv(CELL_PASS_DIR/f'{dataset}_valid.tsv', sep='\t', index=None, na_rep='')
    
    df_t.to_csv(CELL_PASS_DIR/f'{dataset}_train.tsv', sep='\t', index=None, na_rep='')
    print('Source: {:<6} | DB-pass: {:<6} | Train cells: {:<6d} | Valid cells: {:<6d}'.format(dataset, df.shape[0], df_t.shape[0], df_tp.shape[0]))

"""
## Step5.1: Find optimal distance threshold

找到最优的clonotype聚类阈值
"""

# %%
sh_file = SCRIPT / 'find_threshold.sh'
sh_out = SCRIPT / f'log/slurm_find_threshold.out'

with open(sh_file, 'w') as f:
    f.write('#!/bin/bash\n\n')

    f.write('#SBATCH -n 1 # 指定核心数量\n')
    f.write('#SBATCH -N 1 # 指定node的数量\n')
    f.write('#SBATCH -t 0-02:00 # 运行总时间，天数-小时数-分钟, D-HH:MM\n')
    f.write('#SBATCH --mem=5G # 所有核心可以使用的内存池大小, MB为单位\n')
    f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
    f.write('#SBATCH --nodelist=ab03\n\n')

    f.write("docker start {}\n".format(CONTAINER))
    now = datetime.now().strftime('%Y%m%d%H%M%S')
    stdout = WORK_DIR / f'data/log/findThreshold_{now}.out'
    code = f'''docker exec {CONTAINER} /bin/bash -c "(time Rscript findThreshold.r) > {stdout} 2>&1"\n'''
    f.write('{}\n'.format(code))

os.system(f'sbatch {sh_file}')

# %%
"""
## Step5.2: Define clones

进行clonotype聚类
"""
df = pd.read_csv(THRESHOLD_FILE)
for row in df.iloc[66:, :].itertuples(False):
    if pd.isna(row[1]):
        print(f'No threshold of {row[0]}')
        continue
    print(row[0])
    define_clones(Path(row[0]), row[1], CONTAINER, True)

# %%
"""
## Step6: Infer germline sequence

为每个clonotype推断其germline sequence
"""
for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    dbs = (DATA_HOME/'clones').glob(f'{dataset}*')
    for db in dbs:
        db_in_docker = WORK_DIR/'data/clones' / db.name
        create_germline(db_in_docker, CONTAINER, True)

# %%
"""
## Step7: Merge *_aa columns in fmt19 with *clone_pass_germ-pass.tsv

合并fmt19格式的注释信息与fmt7格式的文件，得到完整注释信息的序列文件
"""

for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    # if dataset=='CM00':
    #     germ_f = HOME/f'clones/{dataset}_train-db_pass_clone-pass_germ-pass.tsv'
    # else:
    #     germ_f = HOME/f'clones/{dataset}_valid-db_pass_clone-pass_germ-pass.tsv'
    
    db = DATA_HOME/f'fmt19/{dataset}.fmt19'
    df_fmt19 = pd.read_csv(db, sep='\t')
    cols = [i for i in df_fmt19.columns if i.find('_aa')>=0]
    cols.remove('junction_aa')
    cols.append('sequence_id')
    
    germ_files = (DATA_HOME/'clones').glob(f'{dataset}*germ-pass*')
    for germ_f in germ_files:
        df = pd.read_csv(germ_f, sep='\t')
        df2 = pd.merge(df, df_fmt19.loc[:, cols], on='sequence_id', how='left')
        df2 = fix_wrong_dtype(df2)
        df2.to_csv(germ_f.parent/f'{germ_f.stem}_merged.tsv', sep='\t', index=None, na_rep='')
        print('Source: {:<6} | File: {:<} | Source shape: {:<6d} | Merged shape: {:<6d}'.format(dataset, germ_f.stem, df.shape[0], df2.shape[0]))

# %%
"""
## Last setp: 统计处理完的细胞及序列数量

"""

for dataset in [f'CM{i}' for i in CURRENT_DATASET]:
    clone_f_train = DATA_HOME/f'clones/{dataset}_train-clone_pass.tsv'
    germ_f_train = DATA_HOME/f'clones/{dataset}_train-clone_pass_germ-pass.tsv'
    df_clone_train = pd.read_csv(clone_f_train, sep='\t')
    df_germ_train = pd.read_csv(germ_f_train, sep='\t')
    
    clone_f_valid = DATA_HOME/f'clones/{dataset}_valid-clone_pass.tsv'
    germ_f_valid = DATA_HOME/f'clones/{dataset}_valid-clone_pass_germ-pass.tsv'
    if clone_f_valid.exists():
        df_clone_valid = pd.read_csv(clone_f_valid, sep='\t')
        df_germ_valid = pd.read_csv(germ_f_valid, sep='\t')
        
        statistic[dataset].extend([df_clone_train.shape[0],
                                    df_clone_train.shape[0],
                                    len(set(df_clone_train.sequence_alignment)),
                                   df_germ_train.shape[0],
                                    len(set(df_germ_train.sequence_alignment)),
                                   df_clone_valid.shape[0],
                                    df_clone_valid.shape[0],
                                    len(set(df_clone_valid.sequence_alignment)),
                                   df_germ_valid.shape[0],
                                    len(set(df_germ_valid.sequence_alignment))])
        
    else:
        clone_f_valid_tp = DATA_HOME/f'clones/{dataset}_valid_tp-clone_pass.tsv'
        germ_f_valid_tp = DATA_HOME/f'clones/{dataset}_valid_tp-clone_pass_germ-pass.tsv'
        clone_f_valid_tn = DATA_HOME/f'clones/{dataset}_valid_tn-clone_pass.tsv'
        germ_f_valid_tn = DATA_HOME/f'clones/{dataset}_valid_tn-clone_pass_germ-pass.tsv'
        
        df_clone_valid_tp = pd.read_csv(clone_f_valid_tp, sep='\t')
        df_germ_valid_tp = pd.read_csv(germ_f_valid_tp, sep='\t')
        df_clone_valid_tn = pd.read_csv(clone_f_valid_tn, sep='\t')
        df_germ_valid_tn = pd.read_csv(germ_f_valid_tn, sep='\t')
        
        statistic[dataset].extend([df_clone_train.shape[0],
                                    df_clone_train.shape[0],
                                    len(set(df_clone_train.sequence_alignment)),
                                   df_germ_train.shape[0],
                                    len(set(df_germ_train.sequence_alignment)),
                                   df_clone_valid_tp.shape[0],
                                   df_clone_valid_tp.shape[0],
                                    len(set(df_clone_valid_tp.sequence_alignment)),
                                   df_germ_valid_tp.shape[0],
                                    len(set(df_germ_valid_tp.sequence_alignment)),
                                   df_clone_valid_tn.shape[0],
                                    df_clone_valid_tn.shape[0],
                                    len(set(df_clone_valid_tn.sequence_alignment)),
                                   df_germ_valid_tn.shape[0],
                                    len(set(df_germ_valid_tn.sequence_alignment))])
        
print(statistic)

with open(DATA_HOME/'statistic.csv', 'a+') as f:
    f_csv = csv.writer(f)
    # f_csv.writerow(statics_header)
    f_csv.writerows(statistic.values())
