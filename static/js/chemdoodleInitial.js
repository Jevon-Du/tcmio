/*  @js 所有chemdoodle窗口初始化
***************************************************************************************************************/
var TDMap = getMap();
var canMap = getMap();

function getMap() {
    var map_ = new Object();    
    map_.put = function(key, value) {    
        map_[key+'_'] = value;    
    };    
    map_.get = function(key) {    
        return map_[key+'_'];    
    };    
    map_.remove = function(key) {    
        delete map_[key+'_'];  
    };     
    return map_;
}

function structureInitial(){
    canMap.put(0,sketcher0);
    canMap.put(1,sketcher1);
    canMap.put(2,sketcher2);
    canMap.put(3,sketcher3);
    canMap.put(4,sketcher4);
    // canMap.put(5,sketcher5);
    // canMap.put(6,sketcher6);
    // canMap.put(7,sketcher7);
    // canMap.put(8,sketcher8);
    // canMap.put(9,sketcher9);
    // canMap.put(10,sketcher10);
    // canMap.put(11,sketcher11);
    // canMap.put(12,sketcher12);
    // canMap.put(13,sketcher13);
    // canMap.put(14,sketcher14);
    sketcher0.clear();
    sketcher1.clear();
    sketcher2.clear();
    sketcher3.clear();
    sketcher4.clear();
    // sketcher5.clear();
    // sketcher6.clear();
    // sketcher7.clear();
    // sketcher8.clear();
    // sketcher9.clear();
    // sketcher10.clear();
    // sketcher11.clear();
    // sketcher12.clear();
    // sketcher13.clear();
    // sketcher14.clear();
}

function hoverMolInitial(){
    canMap.put(71,sketcher71);
    canMap.put(72,sketcher72);
    sketcher71.clear();
    sketcher72.clear();
    
}
function selectMolInitial(){
    canMap.put(73,sketcher73);
    canMap.put(74,sketcher74);
    sketcher73.clear();
    sketcher74.clear();
}

//  单个chemdoodle窗口初始化
function chemdoodleInitial(data){
    var molecule = ChemDoodle.readMOL(data);
    ChemDoodle.informatics.removeH(molecule);
    stickTransformer.loadMolecule(molecule);
}

//set default style
function structureStyle(transformer){
    transformer.specs.atoms_useJMOLColors = true;
    // make bonds thicker
    transformer.specs.bonds_width_2D = 1;
    // don't draw atoms
    transformer.specs.atoms_display = true;
    // change the background color to black
    //transformer.specs.backgroundColor = '#fbfcfd';
    transformer.specs.backgroundColor = '#fff';
    // clear overlaps to show z-depth
    transformer.specs.bonds_clearOverlaps_2D = true;
} 