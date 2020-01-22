const tileDimX = 16
const tileDimY = 16
const scale = 3

const roomstr = `width 15
height 20
layer 0
318,319,319,319,319,319,319,319,319,319,319,319,319,319,320
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,134,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
15,16,16,16,16,16,16,16,16,16,16,16,16,16,17
45,46,46,46,46,46,46,46,46,46,46,46,46,46,47
layer 0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,468,0,307,307,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,307,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,255,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 1
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,134,0,0,0,0,0,0,0,0,0
0,0,132,133,133,0,0,0,0,0,132,133,0,0,0
0,0,0,0,407,0,0,0,0,0,0,0,132,134,132
0,0,227,0,0,0,0,0,0,0,0,0,0,0,0
0,132,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 1
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,438,0,0,0,0,0,0,0,0,0
0,0,0,0,407,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 2
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,13,0,13,0,0,0,0,0,0,0,0,0
0,0,13,13,13,0,0,0,440,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,42,44,0
0,288,290,0,0,0,0,0,0,0,0,0,0,0,0
0,42,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 3
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,102,103,103,103,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,102
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 4
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,440,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,440,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 5
0,0,0,0,0,0,0,0,0,0,102,103,103,103,104
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,102,103,103,103,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,72
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 6
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,45
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 7
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,72,73,73,73,74
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
layer 8
0,0,0,0,0,0,0,0,0,0,288,289,289,289,290
0,0,0,0,0,0,0,0,0,0,12,13,13,13,14
0,0,0,0,0,0,0,0,0,0,42,43,43,43,44
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
collider block,x 0,y 0,z -1,w 15,h 20,d 1,name floor
collider block,x -1,y -1,z 0,w 1,h 22,d 10,name leftwall
collider block,x 15,y -1,z 0,w 1,h 22,d 10,name rightwall
collider block,x -1,y -1,z 0,w 17,h 1,d 10,name upwall
collider block,x -1,y 20,z 0,w 17,h 1,d 10,name downwall
collider block,x 1,y 8,z 0,w 2,h 1,d 1,name platform_chunk1
collider block,x 1,y 9,z 0,w 1,h 1,d 1,name platform_chunk2
collider triangle,x 2,y 9,z 0,d 1,rx2 0.96,ry2 0,rx3 0,ry3 0.96,axis z,name platform_chunk3
collider triangle,x 4,y 7,z -0.04,d 1,rx2 0,ry2 1,rx3 1,ry3 0,axis x,name platform_chunk4
collider triangle,x 5,y 6,z -0.04,d 1,rx2 0,ry2 1,rx3 1,ry3 0,axis y,name platform_chunk5
collider block,x 3,y 5,z 0,w 1,h 2,d 1,name block
collider block,x 2,y 6,z 0,w 1,h 1,d 1,name blockl
collider block,x 4,y 6,z 0,w 1,h 1,d 1,name blockr
collider block,x 5,y 5,z 0,w 1,h 1,d 1,name blocktop
collider block,x 5.1,y 9.4,z 0,w 0.8,h 0.45,d 0.5,name rock
collider block,x 7,y 6,z 1.96,w 1,h 1,d 0.1,name platform1
collider block,x 8,y 6,z 0.96,w 1,h 1,d 0.1,name platform2
collider block,x 11,y 7,z 1.96,w 1,h 1,d 0.1,name platform3
collider block,x 10,y 4,z 0,w 5,h 3,d 4,name bigblock
collider block,x 12,y 7,z 0,w 2,h 1,d 1,name bigblocka
collider block,x 14,y 7,z 0,w 1,h 1,d 3,name bigblockb
`

function writeRoom() {
    write(window.room)
}

function write(room) {
    console.log("writing room", { room })
    const w = room.width
    let file = `width ${w}\nheight ${room.height}\n`
    file += room.layers.map(lyr => {
        let join = []
        let tiles = lyr.tiles
        if (tiles.length <= w) {
            join.push(tiles.map(t => t ? t : '0').join(','))
        } else {
            for (let i = 0; i < tiles.length; i++) {
                if (i === tiles.length - 1) {
                    join.push(tiles.slice(i - w + 1, i + 1).map(t => t ? t : '0').join(','))
                } else if (i % w === 0 && i !== 0) {
                    join.push(tiles.slice(i - w, i).map(t => t ? t : '0').join(','))
                }
            }
        }
        return `layer ${lyr.priority}\n${join.join('\n')}`
    }).join('\n')
    // TODO: colliders
    file += room.colliders.map(col => {
        return Object.keys(col)
            .map(k => `${k === 'type' ? 'collider' : k} ${col[k]}`)
            .join(',')
    }).join('\n')
    console.log({ file })
    return file
}

function parse(file) {
    function padArray(arr, len) {
        let n = Array(len - arr.length).fill('0')
        return arr.concat(n)
    }

    function createLayer(tiles, priority) {
        return {
            tiles: tiles.map(t => parseInt(t))
                .concat(Array((room.width * room.height) - tiles.length).fill(0)),
            priority
        }
    }

    const room = { layers: [], colliders: [] }
    let curr = []
    let currPr = 0
    let split = file.split('\n')
    for (let i = 0; i < split.length; i++) {
        if (split[i].trim() === "") {
            continue
        }

        let line = split[i].split(',')
        if (line[0].startsWith("collider")) {
            // TODO
        } else if (line[0].startsWith("width")) {
            let sp = line[0].split(' ')
            room.width = parseInt(sp[1])
        } else if (line[0].startsWith("height")) {
            let sp = line[0].split(' ')
            room.height = parseInt(sp[1])
        } else if (line[0].startsWith("layer")) {
            if (curr.length > 0) {
                room.layers.push(createLayer(curr, currPr))
            }

            curr = []
            currPr = parseInt(line[0].split(' ')[1])
        } else { // case for tiles
            curr = curr.concat(padArray(line, room.width || 20))
        }
    }
    room.layers.push(createLayer(curr, currPr))
    return room
}

function drawRoom(room) {
    let canvas = document.getElementById('editor')
    const roomW = room.width
    const roomH = room.height

    if (canvas.getContext) {
        let ctx = canvas.getContext('2d')
        ctx.webkitImageSmoothingEnabled = false
        ctx.msImageSmoothingEnabled = false
        ctx.imageSmoothingEnabled = false

        ctx.clearRect(0, 0, roomW * scale * tileDimX, roomH * scale * tileDimY)

        for (let i = 0; i < room.layers.length; i++) {
            if (i === selectedLayer()) {
                ctx.globalAlpha = 1.0
            } else {
                ctx.globalAlpha = 0.25
            }
            let tiles = room.layers[i].tiles
            for (let j = 0; j < tiles.length; j++) {
                let tile = tiles[j]
                let sx = (tile % 30) * 16
                let sy = Math.floor(tile / 30) * 16

                ctx.drawImage(document.getElementById('tileset'),
                    sx, sy, tileDimX, tileDimY,
                    ((j % roomW) * 16) * scale, (Math.floor(j / roomW) * 16) * scale,
                    tileDimX * scale, tileDimY * scale)
            }
        }


    }
}

function editorClick(ev) {
    let layerNum = selectedLayer()

    let x = ev.offsetX,
        y = ev.offsetY

    let tileX = Math.floor((x / scale) / tileDimX),
        tileY = Math.floor((y / scale) / tileDimY)

    let num = (tileX % window.room.width) + (tileY * window.room.width)

    window.room.layers[layerNum].tiles[num] = selectedTile()
    drawRoom(window.room)
    // let tile = window.room.layers[layerNum].tiles[num]
    // if (!tile) { tile = 0 }
    // console.log('Tile clicked', { x, y, tileX, tileY, num, tile })
}

function editorRightClick(ev) {
    ev.preventDefault()
    let layerNum = selectedLayer()

    let x = ev.offsetX,
        y = ev.offsetY

    let tileX = Math.floor((x / scale) / tileDimX),
        tileY = Math.floor((y / scale) / tileDimY)

    let num = (tileX % window.room.width) + (tileY * window.room.width)

    let tile = window.room.layers[layerNum].tiles[num]
    if (!tile) tile = 0
    setSelectedTile(tile)
    // let tile = window.room.layers[layerNum].tiles[num]
    // if (!tile) { tile = 0 }
    // console.log('Tile clicked', { x, y, tileX, tileY, num, tile })
}

function loadCanvas(room) {
    let old = document.getElementById('editor')
    if (old) {
        document.removeChild(old)
    }
    let canv = document.createElement('canvas')
    canv.id = 'editor'
    canv.width = room.width * tileDimX * scale
    canv.height = room.height * tileDimY * scale
    canv.addEventListener('click', editorClick)
    canv.addEventListener('contextmenu', editorRightClick)
    document.getElementById('editcontainer').appendChild(canv)
}

function loadPick(picked) {
    const COLS = 30
    let row = Math.floor(picked / COLS),
        col = picked % COLS

    let x = col * tileDimX,
        y = row * tileDimY

    let picker = document.getElementById('tilepicker')
    if (picker.getContext) {
        let ctx = picker.getContext('2d')
        ctx.webkitImageSmoothingEnabled = false
        ctx.msImageSmoothingEnabled = false
        ctx.imageSmoothingEnabled = false

        ctx.clearRect(0, 0, 480, 250)

        ctx.drawImage(document.getElementById('tileset'), 0, 0, 480, 256)

        ctx.strokeStyle = 'yellow'
        ctx.lineWidth = 3
        ctx.strokeRect(x + 0.5, y + 0.5, tileDimX, tileDimY)

        ctx.strokeStyle = 'black'
        ctx.lineWidth = 1
        ctx.strokeRect(x - 1.5, y - 1.5, tileDimX + 3, tileDimY + 3)
        ctx.strokeRect(x + 1.5, y + 1.5, tileDimX - 3, tileDimY - 3)
    }

    let selected = document.getElementById('tileselected')
    if (selected.getContext) {
        let ctx = selected.getContext('2d')
        ctx.webkitImageSmoothingEnabled = false
        ctx.msImageSmoothingEnabled = false
        ctx.imageSmoothingEnabled = false

        ctx.clearRect(0, 0, 16 * scale, 16 * scale)

        ctx.drawImage(document.getElementById('tileset'),
            x, y, tileDimX, tileDimY,
            0, 0, tileDimX * scale, tileDimY * scale)
    }
}

function pickClick(ev) {
    let x = ev.offsetX,
        y = ev.offsetY

    let tileX = Math.floor(x / tileDimX),
        tileY = Math.floor(y / tileDimY)

    let tile = (tileX % 30) + (tileY * 30)
    setSelectedTile(tile)
}

function selectedTile() {
    return window.selected_tile
}
function setSelectedTile(tile) {
    window.selected_tile = tile
    loadPick(tile)
}

function buildLayersList(room) {
    // empty layers
    document.getElementById('layercontainer').innerHTML = ""

    for (let i = 0; i < room.layers.length; i++) {
        let optwrapper = document.createElement('div')
        let opt = document.createElement('input')
        opt.name = "layers"
        opt.type = "radio"
        opt.id = `layer-${i}`
        opt.value = i
        if (i === 0) {
            opt.setAttribute('checked', true)
        }
        optwrapper.appendChild(opt)
        optwrapper.innerHTML += `<label for="layer-${i}">Layer ${i} (priority ${room.layers[i].priority})</label>`
        optwrapper.className = 'layer-option'
        document.getElementById('layercontainer').appendChild(optwrapper)
        document.getElementById(`layer-${i}`).addEventListener('change', onLayerSelect)
    }
}

function selectedLayer() {
    return parseInt(document.querySelector('input[name="layers"]:checked').value)
}
function onLayerSelect() {
    drawRoom(window.room)
}

function start() {
    window.selected_tile = 13
    window.room = parse(roomstr)
    loadCanvas(room)
    buildLayersList(room)
    document.getElementById('tilepicker').addEventListener('click', pickClick)
    loadPick(selectedTile())
    drawRoom(room)
}