
function getCatalogs(){
    return [
        {id: "998", name: "分类1", parent_id: ""},
        {id: "997", name: "分类2", parent_id: ""},
        {id: "486", name: "子分类",parent_id: "997"}
    ]
}


function getCascaderFormat(){
    let catalogs = getCatalogs();
    let tmpMap = {};
    let ret = [];

    catalogs.forEach(function(catalog){
        tmpMap[catalog.id] = {
            "value" : catalog.id,
            "parent_id": catalog.parent_id,
            "label": catalog.name,
            "children": [],
        };
    });

    for (const [id, catalog] of Object.entries(tmpMap)) { // eslint-disable-line no-unused-vars
        if(catalog.parent_id === ""){
            ret.push(catalog);
            continue;
        }
        if (!(Object.prototype.hasOwnProperty.call(tmpMap, catalog.parent_id))){
            continue;
        }
        tmpMap[catalog.parent_id].children.push(catalog)
    }

    for (const [id, catalog] of Object.entries(tmpMap)) { // eslint-disable-line no-unused-vars
        if(catalog.children.length === 0){
            delete catalog.children;
        }
    }

    return ret;
}

export default {
    getCatalogs,
    getCascaderFormat
}