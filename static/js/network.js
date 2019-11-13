/*  @js  vis网络初始化
--------------------------------------------------------------------------------------------------------*/
var nodes = null;
var edges = null;
var network = null;
var highlightActive = false;

function destroy() {
    if (network !== null) {
        network.destroy();
        network = null;
    }
}
/*  Draw Network in result.html*/
function draw(node_data, edge_data) {
    destroy();
    // create a network
    nodes = node_data;
    edges = edge_data;
    var container = document.getElementById('relationChart');
    var data = {
        nodes: nodes,
        edges: edges
    };

    var options = {

        groups: {
            'useDefaultGroups': false,
            'prescriptions': {
                font: {
                    color: "rgba(0,0,0,1)",
                    size: 20, // px
                    face: 'Arimo',
                    background: 'none',
                    strokeWidth: 2, // px
                    strokeColor: '#fff',
                    align: 'center'
                },
                shape: 'triangle',
                color: '#FF9900' // orange
            },
            'tcms': {
                shape: 'square',
                color: "#C5000B" // red
            },
            'ingredients': {
                size: 20,
                shape: 'dot',
                color:'#97C2FC' // blue
            },
            'ligands': {
                size: 20,
                shape: 'star',
                color:'#7BE141'
                // color: "#109618" // purple
            },
            'targets': {
                font: {
                    size: 20, // px
                    face: 'Arimo',
                    background: 'none',
                    strokeWidth: 0, // px
                    strokeColor: '#ffffff',
                    align: 'center'
                },
                shape: 'diamond',
                color: '#FB7E81'
            }
        },
        edges: {
            arrows: {
                middle: {
                    enabled: true,
                    scaleFactor: 2
                }
            },
            arrowStrikethrough: true,
            color: {
                color: '#3399CC',
                highlight: '#CC3366',
                hover: '#CC3366',
                opacity: 0.5
            },
            width: 2.0,
            hoverWidth: 5.0,
            smooth: {
                enabled: true,
                type: "dynamic",
                roundness: 0.5
            }
        },
        layout: {
            randomSeed: 34
        },
        physics: {
            forceAtlas2Based: {
                gravitationalConstant: -20,
                centralGravity: 0.005,
                springLength: 240,
                springConstant: 0.205,
                avoidOverlap: 0.3
            },
            maxVelocity: 148,
            solver: 'forceAtlas2Based',
            timestep: 0.22,
            stabilization: {
                enabled: true,
                iterations: 50,
                updateInterval: 25
            }
        },
        interaction: {
            navigationButtons: true,
            keyboard: true,
            hover: true,
            hideEdgesOnDrag: true
        },
    };
    network = new vis.Network(container, data, options);

    // network.on("stabilizationProgress", function(params) {
    //     var maxWidth = $("#border").width();
    //     var minWidth = 20;
    //     var widthFactor = params.iterations/params.total;
    //     var width = Math.max(minWidth,maxWidth * widthFactor);
    //     document.getElementById('bar').style.width = width + 'px';
    //     document.getElementById('text').innerHTML = Math.round(widthFactor*100) + '%';
    // });
    // network.once("stabilizationIterationsDone", function() {
    //     document.getElementById('text').innerHTML = '100%';
    //     document.getElementById('bar').style.width = '85%';
    //     document.getElementById('loadingBar').style.opacity = 0;
    //     // really clean the dom element
    //     setTimeout(function () {document.getElementById('loadingBar').style.display = 'none';}, 500);
    // });

    network.on("selectNode", function (){
        //$("#mask").show();
        //$(".mol-detail").show();

        // 1. 获取当前节点的id --> ajax的参数
        // 2. 获取当前节点的group --> ajax的url
        var selectNodeID = this.getSelectedNodes()[0];
 
        var selectedNodeGroup = null;
        var nodesCount = node_data.length;
        for (var i = 0; i < nodesCount; i++) {
            if (selectNodeID == node_data[i].id){
                selectedNodeGroup = node_data[i].group;
                break;
            }
        }

        console.log(selectNodeID);
        console.log(selectedNodeGroup);

        // TODO: 这里的数据跟后端绑定了, 需要解耦
        if (selectedNodeGroup=='tcms'){
            selectNodeID = selectNodeID - 1000;
        } else if (selectedNodeGroup=='prescriptions'){
            selectNodeID = selectNodeID - 10000;
        } else if (selectedNodeGroup=='ingredients'){
            selectNodeID = selectNodeID - 100000;
        } else if (selectedNodeGroup=='ligands'){
            selectNodeID = selectNodeID - 200000;
        }

        // 请求数据
        $.ajax({
            url: selectedNodeGroup + '/'  + selectNodeID +'/json',
            type: "get",
            dataType: "json",
            error: function (jqXHR, textStatus, errorThrown) {
                if (textStatus == "timeout") {
                    alert("Request timeout, please refresh the page.");
                } else {
                    alert(textStatus);
                }
            }
        }).done(function (msg) {
            console.log(msg);

            render(msg.Data, selectedNodeGroup);

            if (selectedNodeGroup=='ligands' || selectedNodeGroup=='ingredients'){
                show_mol_structure(msg.Data.Mol, selectedNodeGroup);
            }

            $(".modal").modal("show");
        });
    });
}

function which_type(data_type){
    var name=capitalizeFirstLetter(data_type.slice(0,-1));
    return name;
}

function render(data, data_type){
    // 先清空数据
    $('.info tbody').empty();

    var title = which_type(data_type);
    if (title=='Prescription'){
        var heading =  title + ' : ' + data.PinyinName;
        $('.item_name').html(heading);
    } else if (title=='Tcm'){
        var heading =  title + ' : ' + data.EnglishName;
        $('.item_name').html(heading);
    } else {
        var name = ('Name' in data ? data.Name : data.Id);
        var heading = which_type(data_type) + ' : ' + name;
        $('.item_name').html(heading);
    }
    var flag = title=='Target' ? true : false
    // 创建表格,渲染数据
    var html = '';
    for (var key in data){
        html += '<tr class="item">';
        html += '<th class="item_name">' + key +'</th>';
        if (data[key]){
            if(key=='Mol'){
                html += '<td class="item_info">' + canvas_format(data_type, 1) +'</td>';
            } else{
                html += '<td class="item_info">' + link_format(data[key], key, flag) +'</td>';
            }
        } else {
            html += '<td class="item_info"><i>N/A<i></td>';
        }
        html += '</tr>';
    };

    $('.info tbody').append(html);
}

//  save img 功能：将canvase先转化为base64字符串，再转为png/jpg/gif/bmp
$(".js-saveCanvas").on("click",function (){
    var mydate = new Date();      
    var year=mydate.getFullYear();
    var month=mydate.getMonth() + 1;
    var day=mydate.getDate();
    var time = year + month + day;

    var oCanvas = document.getElementsByTagName("canvas")[0];
    // 图片导出为 png 格式  
    // var type = $(".btn-xs").attr("name"); 
    var type = 'png'; 
    var imgData = oCanvas.toDataURL(type);

    /** 
     * 获取mimeType 
     * @param  {String} type the old mime-type 
     * @return the new mime-type 
     */  
    var _fixType = function(type) {  
        type = type.toLowerCase().replace(/jpg/i, 'jpeg');  
        var r = type.match(/png|jpeg|bmp|gif/)[0];  
        return 'image/' + r;  
    };  
        
    // 加工image data，替换mime type  
    imgData = imgData.replace(_fixType(type),'image/octet-stream');  

    /** 
     * 在本地进行文件保存 
     * @param  {String} data     要保存到本地的图片数据 
     * @param  {String} filename 文件名 
     */  
    function saveFile(data, filename){  
        var save_link = document.createElementNS('http://www.w3.org/1999/xhtml', 'a');  
        save_link.href = data;  
        save_link.download = filename;  
        
        var event = document.createEvent('MouseEvents');  
        event.initMouseEvent('click', true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);  
        save_link.dispatchEvent(event);  
    };  
        
    // 下载后的问题名  
    var filename = time + 'network.' + type; 
    //console.log(oCanvas.toDataURL("image/png"));
    saveFile(imgData,filename);
});