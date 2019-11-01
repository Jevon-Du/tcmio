/**********************format every link*******************/
function link_format(id_val, type, is_target){
    var db_link = {
        'UniprotId': 'https://www.uniprot.org/uniprot/',
        'ChemblId': 'https://www.ebi.ac.uk/chembl/target_report_card/',
        'EcNumber': 'https://enzyme.expasy.org/EC/',
        'Kegg': 'https://www.genome.jp/dbget-bin/www_bget?',
        'Pdb': 'https://www.rcsb.org/structure/',
        'Drug': 'https://www.drugbank.ca/drugs/'
    }
    if (!is_target){
        db_link['ChemblId'] = 'https://www.ebi.ac.uk/chembl/compound_report_card/'
    }
    var link_icon_el = '&nbsp<span class="glyphicon glyphicon-new-window" aria-hidden="true"></span>';
    if (db_link.hasOwnProperty(type)){
        var new_el = '';
        if (type=='UniprotId'){
            new_el = '<a href="https://www.uniprot.org/uniprot/'+ id_val +'" target="_blank">' + id_val + link_icon_el +'</a>';
        } else {
            var ids = id_val.split(';');
            for (var i in ids){
                new_el = new_el + '<a href="' + db_link[type] + ids[i].replace(' ', '') +'" target="_blank">' + ids[i].replace(' ', '') + link_icon_el +'</a>&nbsp&nbsp';
            }
        }
        return new_el;
    }
    return id_val;
}

function canvas_format(idx){
    var var_name = 'sketcher' + idx;
    var chemdoodle_html = '<canvas id="'+ var_name +'"></canvas>';
    return chemdoodle_html;
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}


function ipm_fromat(idx){
    var iframe_el = '<iframe id="iframe'+ idx +'" name="' + idx +'" src="static/ipmDraw/editor.html" frameborder="0" width="300px" height="300px"></iframe>';
    return iframe_el;
}


function loadMolecule(index, mol) {
    var tempCanvas = canMap.get(index);
    var molecule = ChemDoodle.readMOL(mol);
    ChemDoodle.informatics.removeH(molecule);
    structureStyle(tempCanvas);
    tempCanvas.loadMolecule(molecule);
}

function init_chemdoodle(data_type){
    var interval1 = setInterval(function() {
        if (data_type=='ingredients' || data_type=='ligands'){
            var tds = $('#'+data_type+ ' tbody tr td canvas');
            if (tds) {
                tds.each(function(index,element){
                    var myCanvas = new ChemDoodle.ViewerCanvas('sketcher'+index, 200, 200);
                    myCanvas.emptyMessage = 'No Data Loaded!';
                    myCanvas.repaint();
                    var mol = ChemDoodle.readMOL(MOLS[index]);
                    myCanvas.loadMolecule(mol);
                });
            }
            clearInterval(interval1);
        }
    }, 100);

    var interval2 = setInterval(function() {
        $('.dataTables_length select').addClass('form-control select select-default');
        $('.dataTables_length select').select2({
            dropdownCssClass: 'dropdown-inverse',
            width: '70px',
            allowClear: false
        });
        if ($('.select-default')) {
            clearInterval(interval2);
        }
    }, 100);
}

