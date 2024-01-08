import os
import subprocess
from pathlib import Path
from datetime import datetime
from typing import List, Set, Dict, Tuple, Optional, Union, NoReturn

import pandas as pd

DTYPE = {
    'v_sequence_start': 'Int64', 
    'v_sequence_end': 'Int64', 
    'v_germline_start': 'Int64', 
    'v_germline_end': 'Int64', 
    'v_alignment_start': 'Int64',
    'v_alignment_end': 'Int64',
    'd_sequence_start': 'Int64', 
    'd_sequence_end': 'Int64', 
    'd_germline_start': 'Int64', 
    'd_germline_end': 'Int64', 
    'd_alignment_start': 'Int64',
    'd_alignment_end': 'Int64',
    'j_sequence_start': 'Int64', 
    'j_sequence_end': 'Int64', 
    'j_germline_start': 'Int64', 
    'j_germline_end': 'Int64', 
    'j_alignment_start': 'Int64',
    'j_alignment_end': 'Int64',
    'np1_length': 'Int64', 
    'np2_length': 'Int64',
    'junction_length': 'Int64',
    'junction_aa_length': 'Int64',
    'fwr1_start': 'Int64',
    'fwr1_end': 'Int64',
    'cdr1_start': 'Int64',       
    'cdr1_end': 'Int64',
    'fwr2_start': 'Int64',
    'fwr2_end': 'Int64',
    'cdr2_start': 'Int64',
    'cdr2_end': 'Int64',
    'fwr3_start': 'Int64',
    'fwr3_end': 'Int64',
    'fwr4_start': 'Int64',
    'fwr4_end': 'Int64',
    'cdr3_start': 'Int64',
    'cdr3_end': 'Int64' 
    # 'v_frameshift': str
}

BTYPE = {
    'Memory B-cells': 1,
    'Plasma cells': 2,
    'naive B-cells': 3,
    'Class-switched memory B-cells':4
}

VALID_MICE_old = {
    'CM001': ['VMM_1_1', 'VMM_4_1'],
    'CM002': ['VMM_1', 'VMM_5'],
    'CM003': ['VMM_1', 'VMM_5_9'],
    'CM004': ['VMM_2_1', 'VMM_1_9'],
    'CM005': ['VMM_2_1'],
    'CM006': ['VMM_2_1'],
    'CM008': ['VMM_2_1'],
    'CM009': ['CM009-2-1'],
    'CM010': ['2_1']
}

VALID_MICE = {
    'CM1': ['CM001-1-1', 'CM001-4-1'],
    'CM2': ['CM002-1', 'CM002-5'],
    'CM2_new': ['CM002-2-2-1'],
    'CM3': ['CM003-1', 'CM003-5-9'],
    'CM4': ['CM004-1-9', 'CM004-2-1'],
    'CM5-2': ['2--1'],
    'CM6': ['CM006-2-1'],
    'CM7': ['CM007-2-1'],
    'CM8': ['CM008-2-1'],
    'CM9': ['CM009-2-1'],
    'CM10': ['2_1'],
    'CM11': ['CM-011-AHP-2-1'],
    'CM12': ['1--1'],
    'CM13': ['CM013-BSA-2--1'],
    'CM15': ['CM015-2-1'],
    'CM20': ['CM020-2-1'],
    'CM21': ['CM021-2-1'],
    'CM22': ['CM022-2-1'],
    'CM24': ['CM024-2-1'],
    'CM26': ['CM026-2-1'],
    'CM28': ['CM028-2-1'],
    'CM29': ['CM029-2-1'],
    'CM32': ['CM032-2-1'],
    'CM33': ['CM033-2-1'],
    'CM34': ['CM034-2-1'],
    'CM35': ['CM035-2-1'],
    'CM36': ['CM036-2-1'],
    'CM37': ['CM037-2-1'],
    'CM38': ['CM038-2-1'],
    'CM39': ['CM039-2-1'],
    'CM41': ['CM041-2-1'],
    'CM4_new': ['CM004-2-2-1'],
    'CM43': ['CM043-2-1'],
    'CM6_new': ['CM006-2-2-1'],
    'CM8_new': ['CM008-2-1'],
    'CM40': ['CM-040-2-1'],
    'CM7_new': ['CM-007-2-2-1']
}

def float2int(x: Union[str, float, int]
             ) -> Optional[int]:
    """
    """
    if pd.notna(x) and x:
        if x!='NA' or x!='nan':
            return int(x)
        else:
            return None
    else:
        return None

def fix_wrong_dtype(df: pd.DataFrame) -> pd.DataFrame:
    """
    """
    
    for key, val in DTYPE.items():
        if key in df.columns:
            df[key] = df[key].apply(float2int).astype('Int64')
    return df

def run_igblast(fasta: Path,
                fmt: int = 7,
                use_slurm: bool = True,
                threds: int = 10,
                script: Path = Path('/data/personal/jiewen.du/CM_new/script'),
                blast_path: Path = Path('/data/personal/jiewen.du/ncbi-igblast-1.17.0')
               ) -> NoReturn:
    """
    Run igblastn(v1.17.0).
    
    Args:
        fasta (Path): default is '/data/personal/jiewen.du/CM/data/cell_pass'
        script (Path): default is '/data/personal/jiewen.du/CM/script'
        fmt (int) : outformat of igblastn
        threds (int): Number of threads to run igblastn
        blast_path (str): igblast exec path.
    
    Returns:
        None
    """
    fname = fasta.stem
    ofile = fasta.parent.parent / f'fmt{fmt}/{fname}.fmt{fmt}'
    sh_file = script / f'igblast_{fname}.sh'
    sh_out = script / f'log/igblast_{fname}_fmt{fmt}.out'
    
    with open(sh_file, 'w') as f:
        f.write('#!/bin/bash\n\n')
        
        f.write(f'#SBATCH -n {threds} # 指定核心数量\n')
        f.write('#SBATCH -N 1 # 指定node的数量\n')
        f.write('#SBATCH -t 0-01:00 # 运行总时间，天数-小时数-分钟, D-HH:MM\n')
        f.write('#SBATCH --mem=5G # 所有核心可以使用的内存池大小, MB为单位\n')
        f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
        f.write('#SBATCH --nodelist=ab03\n\n')
        
        if fmt==7:
            code = f'''cd {blast_path}; bin/igblastn -germline_db_V database/imgt_human_ig_v -germline_db_J database/imgt_human_ig_j -germline_db_D database/imgt_human_ig_d -auxiliary_data optional_file/human_gl.aux -ig_seqtype Ig -organism human -show_translation -strand plus -num_threads {threds} -outfmt '7 std qseq sseq btop' -query {fasta} -out {ofile}'''
        else:
            code = f'''cd {blast_path}; bin/igblastn -germline_db_V database/imgt_human_ig_v -germline_db_J database/imgt_human_ig_j -germline_db_D database/imgt_human_ig_d -auxiliary_data optional_file/human_gl.aux -ig_seqtype Ig -organism human -strand plus -num_threads {threds} -outfmt {fmt} -query {fasta} -out {ofile}'''
        
        f.write('{}\n'.format(code))
    
    if use_slurm:
        os.system(f'sbatch {sh_file}')
    else:
        os.system(f'bash {sh_file}')
    
def make_db(db: Path, 
            fasta: Path,
            container: str,
            use_slurm: bool = True,
            ref_db: Path = Path('/data/IMGT'),
            script: Path = Path('/data/personal/jiewen.du/CM/script')
           ) -> NoReturn:
    """
    Create tab-delimited database file to store sequence alignment information.
    
    Args:
        db (Path): IgBLAST output file in format 7 with query sequence, path in docker container.
        fasta (Path): input FASTA files, path in docker container.
        container (str): docker container name or id.
        use_slurm (boool): use slurm to run script.
        ref_db (Path): reference geremline sequence, path in docker container.
        script (Path): Path('/data/personal/jiewen.du/CM/script')
    
    Output file:
        db-pass: database of alignment records with functionality information, V and J calls, and a junction region.
        db-fail: database with records that fail due to no productivity information, no gene V assignment, no J assignment, or no junction region. 
        stdout: stdout and stderr of Makedb.py.
        slurm_out: slurm log file.
    """
    
    ofile = db.parent.parent / f'{db.stem}-db_pass.tsv'
    stdout = db.parent.parent / f'log/{db.stem}-db_pass.out'
    
    vdb = ref_db / 'imgt_human_IGHV.fasta'
    jdb = ref_db / 'imgt_human_IGHJ.fasta'
    ddb = ref_db / 'imgt_human_IGHD.fasta'
    
    sh_file = script / f'make_db_{db.stem}.sh'
    sh_out = script / f'log/slurm_make_db_{db.stem}.out'
    
    with open(sh_file, 'w') as f:
        f.write("#!/bin/bash\n\n")
        f.write('#SBATCH -n 1 # 指定核心数量\n')
        f.write('#SBATCH -N 1 # 指定node的数量\n')
        f.write('#SBATCH -t 0-00:30 # 运行总时间，天数-小时数-分钟, D-HH:MM\n')
        f.write('#SBATCH --mem=5G # 所有核心可以使用的内存池大小, MB为单位\n')
        f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
        f.write('#SBATCH --nodelist=ab03\n\n')
        
        f.write("docker start {}\n".format(container))
       
        code = f'''docker exec {container} /bin/bash -c "MakeDb.py igblast -i {db} -s {fasta} -r {vdb} {jdb} {ddb} -o {ofile} --extended > {stdout} 2>&1"\n'''
        
        f.write('{}\n'.format(code))
    
    if use_slurm:
        os.system(f'sbatch {sh_file}')
    else:
        os.system(f'bash {sh_file}')

def find_threshold(script: Path = Path('/data/personal/jiewen.du/CM/script')):
    sh_file = script / 'find_threshold.sh'
    sh_out = script / f'log/slurm_find_threshold_.out'

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
        stdout = WORK_DIR / LOG_DIR / f'findThreshold_{now}.out'
        code = f'''docker exec {CONTAINER} /bin/bash -c "(time Rscript findThreshold.r) > {stdout} 2>&1"\n'''
        f.write('{}\n'.format(code))
    
    os.system(f'slurm {sh_file}')
    
def define_clones(db: Path, 
                  threshold: float,
                  container: str,
                  use_slurm: bool = True,
                  script: Path = Path('/data/personal/jiewen.du/CM_new/script')
                 ) -> NoReturn:
    """
    
    Args:
       db (Path): tab-delimited file contain sequence alignment information.
    
    Output file:
        
        
    """
    fname = str(db.stem).split('-')[0]
    ofile = db.parent.parent / f'clones/{fname}-clone_pass.tsv'
    log = db.parent.parent / f'clones/log/{fname}_defineClones.log'
    stdout = db.parent.parent / f'clones/log/{fname}_defineClones.out'
    
    sh_file = script / f'define_clones_{fname}.sh'
    sh_out = script / f'log/slurm_define_clones_{fname}.out'
    
    with open(sh_file, 'w') as f:
        f.write('#!/bin/bash\n\n')
        
        f.write('#SBATCH -n 1 # 指定核心数量\n')
        f.write('#SBATCH -N 1 # 指定node的数量\n')
        f.write('#SBATCH -t 0-02:00 # 运行总时间，天数-小时数-分钟, D-HH:MM\n')
        f.write('#SBATCH --mem=5G # 所有核心可以使用的内存池大小, MB为单位\n')
        f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
        f.write('#SBATCH --nodelist=ab03\n\n')
        
        f.write('docker start {}\n'.format(container))

        code = f'''docker exec {container} /bin/bash -c "(time DefineClones.py -d {db} -o {ofile} --act set --model ham --norm len --dist {threshold} --log {log}) > {stdout} 2>&1"\n'''
        
        f.write('{}\n'.format(code))
    
    if use_slurm:
        os.system(f'sbatch {sh_file}')
    else:
        os.system(f'bash {sh_file}')

def create_germline(clone_file: Path,
                    container: str,
                    use_slurm: bool = True,
                    ref_db: Path = Path('/data/IMGT'),
                    script: Path = Path('/data/personal/jiewen.du/CM_new/script')
                   ) -> NoReturn:
    """
    
    Args:
    
    Output file:
        
    """
    fname = clone_file.stem
    log = clone_file.parent / f'log/{fname}_createGermlines.log'
    stdout = clone_file.parent / f'log/{fname}_createGermlines.out'
    
    vdb = ref_db / 'imgt_human_IGHV.fasta'
    jdb = ref_db / 'imgt_human_IGHJ.fasta'
    ddb = ref_db / 'imgt_human_IGHD.fasta'
    
    sh_file = script / f'create_germline_{fname}.slurm'
    sh_out = script / f'log/slurm_create_germline_{fname}.out'
    
    with open(sh_file, 'w') as f:
        f.write('#!/bin/bash\n\n')
        
        f.write('#SBATCH -n 1 # 指定核心数量\n')
        f.write('#SBATCH -N 1 # 指定node的数量\n')
        f.write('#SBATCH -t 0-00:30 # 运行总时间，天数-小时数-分钟, D-HH:MM\n')
        f.write('#SBATCH --mem=5G # 所有核心可以使用的内存池大小, MB为单位\n')
        f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
        f.write('#SBATCH --nodelist=ab03\n\n')
        
        f.write('docker start {}\n'.format(container))
        
        code = f'''docker exec {container} /bin/bash -c "(time CreateGermlines.py -d {clone_file} -g dmask --cloned -r {vdb} {ddb} {jdb} --log {log}) > {stdout} 2>&1"\n'''

        f.write('{}\n'.format(code))

    if use_slurm:
        os.system(f'sbatch {sh_file}')
    else:
        os.system(f'bash {sh_file}')
        
def build_lineage_tree(germ_file: Path,
                       container: str,
                       clone_id: Union[int, str],
                       clone_size: int,
                       min_size: int = 3,
                       remove_cdr3: bool = True,
                       use_slurm: bool = True,
                       script: Path = Path('/data/personal/jiewen.du/CM_new/script')
                      ) -> NoReturn:
    """
    Builds lineage tree and estimates parametres using IgPhyML
    
    Args:
    
    Output files:
    
    """
    
    fname = str(germ_file.stem).split("-")[0]
    dataset = fname.split('_')[0]
    outname = f'{fname}_clone{clone_id}_s{clone_size}'
    outdir = germ_file.parent.parent / 'trees'
    if not outdir.exists():
        outdir.mkdir(parents=True)
    log = outdir / f'log/{outname}_BuildTrees.log'
    stdout = outdir / f'log/{outname}_BuildTrees.out'
    
    sh_file = script / f'build_tree_{fname}_{clone_id}_{clone_size}.sh'
    sh_out = script / f'log/slurm_build_tree_{fname}_{clone_id}_{clone_size}.out'
    
    with open(sh_file, 'w') as f:
        f.write('#!/bin/bash\n\n')
        
        f.write('#SBATCH -n 1 # 指定核心数量\n')
        f.write('#SBATCH -N 1 # 指定node的数量\n')
        
        if (clone_size > 100):
            f.write('#SBATCH --mem=10000MB # 所有核心可以使用的内存池大小，MB为单位\n')
            f.write('#SBATCH -t 0-24:00 # 运行总时间，天数-小时数-分钟， D-HH:MM\n')
        else:
            f.write('#SBATCH --mem=2000MB # 所有核心可以使用的内存池大小，MB为单位\n')
            f.write('#SBATCH -t 0-05:00 # 运行总时间，天数-小时数-分钟， D-HH:MM\n')
        
        f.write(f'#SBATCH -o {sh_out} # 把输出结果STDOUT保存在哪一个文件\n')
        f.write('#SBATCH --nodelist=ab03\n\n')
        
        f.write('docker start {}\n\n'.format(container))
        
        if remove_cdr3:
            code = f'''docker exec {container} /bin/bash -c "(time BuildTrees.py -d {germ_file} --clones {clone_id} --collapse --ncdr3 --igphyml --clean all --minseq {min_size} --outname {outname} --outdir {outdir/dataset} --log {log} --optimize tlr) > {stdout}"\n'''
        else:
            code = f'''docker exec {container} /bin/bash -c "(time BuildTrees.py -d {germ_file} --clones {clone_id} --collapse --igphyml --clean all --minseq {min_size} --outname {outname} --outdir {outdir/dataset} --log {log} --optimize tlr) > {stdout}"\n'''
        
        f.write('{}\n'.format(code))
    
    if use_slurm:
        os.system(f'sbatch {sh_file}')
    else:
        os.system(f'bash {sh_file}')

def build_repertoire_tree(germ_file: Path,
                          container: str,
                          min_size: int = 3,
                          max_size: int = 500,
                          ncdr3: bool = True,
                          use_slurm: bool = True
                         ) -> NoReturn:
    """
    
    Args:
        
    Returns:
    
    """
    df = pd.read_csv(germ_file, sep='\t')
    df.drop_duplicates(subset=['sequence_alignment'], inplace=True)
    clone_size = df['clone_id'].value_counts()
    
    clone_size_df = clone_size[(clone_size>=min_size)&(clone_size<=max_size)]
    
    # 去重后clonetype_size<=500
    for clone_id, size in clone_size_df.iteritems():
        ifile = Path(f'/data/data/clones/{germ_file.name}')
        build_lineage_tree(ifile, container, clone_id, size, min_size, ncdr3, use_slurm)

def rebuild_tree(germ_file: Path,
                 container: str,
                 min_size: int = 3,
                 max_size: int = 500,
                 do_run: bool = True,
                 ncdr3: bool = True,
                 use_slurm: bool = True,
                ) -> NoReturn:
    """
    
    """
    
    df = pd.read_csv(germ_file, sep='\t')
    df.drop_duplicates(subset=['sequence_alignment'], inplace=True)
    clone_size = df['clone_id'].value_counts()
    clone_size_df = clone_size[(clone_size>=min_size)&(clone_size<=max_size)]
    
    fname = str(germ_file.stem).split("-")[0]
    dataset = fname.split('_')[0]
    
    print(f'Rebuild Tree: {germ_file.name}')
    for clone_id, size in clone_size_df.iteritems():
        f_pass = germ_file.parent.parent/f'trees/{dataset}/{fname}_clone{clone_id}_s{size}_igphyml-pass.tab'
        f_fail = germ_file.parent.parent/f'trees/{dataset}/{fname}_clone{clone_id}_s{size}_lineages.tsv'
        if not (f_pass.exists() or f_fail.exists()):
            print(f'Wait to rebiuld. Clone_id: {clone_id} | clone size: {size}')
            if do_run:
                ifile = Path(f'/data/data/clones/{germ_file.name}')
                build_lineage_tree(ifile, container, clone_id, size, min_size, ncdr3, use_slurm)
            