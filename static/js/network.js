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
            0: {
                font: {
                    size: 20, // px
                    face: 'Arimo',
                    background: 'none',
                    strokeWidth: 0, // px
                    strokeColor: '#ffffff',
                    align: 'center'
                },
                size: 30,
                shape: 'square',
                color: 'lime',
                border: 2,
            },
            1: {
                font: {
                    color: "rgba(0,0,0,1)",
                    size: 20, // px
                    face: 'Arimo',
                    background: 'none',
                    strokeWidth: 2, // px
                    strokeColor: '#fff',
                    align: 'center'
                },
                size: 30,
                // shape: 'square',
                shape: 'triangle',
                /*color: 'rgba(238,99,99,1)',*/
                // color: {background:'#F03967', border:'#713E7F'}
                color: 'rgba(251,0,2,1)'
            },
            2: {
                color: 'rgba(255,160,122,1)'
            },
            3: {
                size: 20,
                shape: 'dot',
                color: 'rgba(92,172,238,1)'
                    /*color:'rgba(255,255,0,1)'*/
            },
            4: {
                size: 20,
                shape: 'diamond',
                color: { background: 'pink', border: 'purple' }
                /*color:'rgba(92,172,238,1)'*/
            },
            9: {
                color: 'rgba(255,160,122,1)'
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
}