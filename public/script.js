const grid = document.querySelector("#fieldGame"); // поле для игры
const letterArray = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'];
let selectedShip = { // выбранный тип корабля
    'line': '',
    'size': ''
};
let ship = [];	// выделенный на поле корабль

// game field
const drawGrid = () => {
    for (var i = 0; i < 100; i++) {
        const cell = document.createElement('div');
        const char = letterArray[i % 10];
        const num = Math.trunc((i) / 10) + 1;

        cell.classList.add("cell");
        cell.setAttribute("id", char + num);
        cell.setAttribute("onmouseover", "cellOnmouseOver(event)");
        cell.setAttribute("onmouseout", "cellOnmouseOut()");
        cell.setAttribute("onclick", "cellOnclick()");

        grid.appendChild(cell);
    }
};
drawGrid();

const getAllCells = () => {
    const cellsList = document.querySelectorAll(".cell");
    let allCells = [];
    for (var i = 0; i < cellsList.length; i++) {
        allCells.push(cellsList[i].getAttribute('id'));
    }
    return allCells;
}

// выбор типа корабля
const choiceShip = (element) => {
    selectedShip.line = element.dataset.line;
    selectedShip.size = Number(element.dataset.size);
}

// наведение курсора на ячейку поля
const cellOnmouseOver = (e) => {
    const thisCell = e.target.getAttribute('id');  // клетка под курсором
    const coordinates = {  // координаты клетки под курсором
        'char': thisCell.charAt(0),
        'num': Number(thisCell.substr(1))
    };
    let allCells = getAllCells();
    let thisShipCills = [];  // массив клеток "под кораблем"

    if (selectedShip.line === 'vertically') {
        allCells.map(cell => {
            const cellChar = letterArray.indexOf(cell.charAt(0));
            if (Number(cell.substr(1)) === coordinates.num) {
                if (letterArray.indexOf(coordinates.char) - selectedShip.size < cellChar
                    && cellChar <= letterArray.indexOf(coordinates.char)) {
                    thisShipCills.push(cell);
                }
            }
        });
    } else if (selectedShip.line === 'horizontally') {
        allCells.map(cell => {
            const cellNum = Number(cell.substr(1));
            if (cell.charAt(0) === coordinates.char) {
                if (coordinates.num - selectedShip.size < cellNum && cellNum <= coordinates.num) {
                    thisShipCills.push(cell);
                }
            }
        });
    }
    thisShipCills.map(cell => {
        document.querySelector('#' + cell).classList.add("cellHover");
    });
    ship = thisShipCills;
};

// уход курсора с ячейки поля
const cellOnmouseOut = () => {
    let allCells = getAllCells();
    allCells.map(cell => {
        const a = document.querySelector('#' + cell);
        if (a.classList.contains("cellHover")) {
            a.classList.remove("cellHover")
        }
    })
}

// установка корабля на поле
const cellOnclick = () => {
    ship.map(cell => {
        document.querySelector('#' + cell).classList.add('cell2');
    })
};

// get layout of the ships as a 2D array
const getLayout = () => {
    let num = 10;
    let char = letterArray;
    var card = [];
    for (var i = 0; i < char.length; i++) {
        card[i] = [];
        for (var j = 0; j < num; j++) {
            const a = document.querySelector('#' + char[i] + (j + 1));
            card[i][j] = a.classList.contains('cell2') ? 1 : 0;
        }
    }
    return card;
};


const sendJSON = (obj) => {
    (async () => {
        const rawResponse = await fetch('/prove', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(obj)
        });
        console.log(JSON.stringify(obj));
        const content = await rawResponse.json();
        console.log(content);
    })();
};

window.onload = () => {
    document.getElementById("sendLayoutBtn").addEventListener("click", sendLayout);
};

const sendLayout = () => {
    const layout = getLayout();
    sendJSON(layout);
};


