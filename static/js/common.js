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
    if (db_link.hasOwnProperty(type)){
        var new_el = '';
        if (type=='UniprotId'){
            new_el = '<a href="https://www.uniprot.org/uniprot/'+ id_val +'" target="_blank">' + id_val +'</a>';
        } else {
            var ids = id_val.split(';');
            for (var i in ids){
                new_el = new_el + '<a href="' + db_link[type] + ids[i].replace(' ', '') +'" target="_blank">' + ids[i].replace(' ', '') + '</a>&nbsp&nbsp';
            }
        }
        return new_el;
    }
    return id_val;
}

function canvas_format(type, idx){
    var id_name = type + '_sketcher' + idx;
    var chemdoodle_html = '<canvas id="'+ id_name +'"></canvas>';
    return chemdoodle_html;
}

// TODO: msg.Data or msg.data
function data_format(msg, type){
    if (type=='targets'){
        var col_name = ['Name', 'GeneName', 'Function', 'ProteinFamily', 'UniprotId', 'ChemblId', 'EcNumber'];
        msg.data.forEach(function(currentValue, index){
            for (var idx in col_name) {
                var new_val = link_format(currentValue[col_name[idx]], col_name[idx], true);
                msg.data[index][col_name[idx]] = new_val;
            }
        });
    } else if(type=='ingredients' || type=='ligands'){
        MOLS.length = 0;
        msg.data.forEach(function(currentValue, index){
            MOLS.push(msg.data[index]['Mol']);
            msg.data[index]['Mol'] = canvas_format(type, index);
        });
    }
    msg.data.forEach(function(currentValue, index){
        msg.data[index]['Id'] = '<a href="/' + type +'/' + currentValue['Id'] + '">' + currentValue['Id'] +'</a>';
    });

    return msg.data;
};

function analyze_data_format(msg, type){
    if (type=='ingredients'){
        msg.data.forEach(function(currentValue, index){
            var ids = currentValue['IngredientID'].split(';');
            var new_el = '';
            for (var i in ids){
                new_el = new_el + '<a href="/ingredients/' + currentValue['IngredientID'] + '" target="_blank">' + currentValue['IngredientID'] +'</a>&nbsp&nbsp';
            };
            msg.data[index]['IngredientID'] = new_el;
        });
    } else if(type=='ligands'){
        msg.data.forEach(function(currentValue, index){
            msg.data[index]['TargetChemblId'] = link_format(currentValue['TargetChemblId'], 'ChemblId', true);
            msg.data[index]['MolChemblId'] = link_format(currentValue['MolChemblId'], 'ChemblId', false);
        });
    }

    return msg.data;
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
                    var myCanvas = new ChemDoodle.ViewerCanvas(data_type+'_sketcher'+index, 200, 200);
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
        $('.dataTables_length select').addClass('form-control select select-primary');
        $('.dataTables_length select').select2({
            // dropdownCssClass: 'dropdown-inverse',
            width: '70px',
            allowClear: false
        });

        if ($('.select-default')) {
            clearInterval(interval2);
        }
    }, 100);
}

