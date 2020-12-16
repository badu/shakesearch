const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    }).catch(err=>{
        console.error(`there was an error calling the server : ${err.message}`)
    });
  },

  updateTable: (results) => {
    const shake = document.getElementById("shake");
    if (shake !== undefined){
        shake.parentNode.removeChild(shake);
    }
    const table = document.getElementById("table-body");
    // remove previous results
    var rowCount = table.rows.length;
    for (var x=rowCount-1; x>0; x--) {
        table.deleteRow(x);
    }

    const rows = [];
    if (results == null){
        rows.push(`<tr>No results found.</tr>`)
    }else{
        for (let result of results) {
          rows.push(`<tr>${result}</tr>`);
        }
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
