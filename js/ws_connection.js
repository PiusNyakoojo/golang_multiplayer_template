// Setup environment
var addr = "localhost:8081"; // when selected: wss -> ws
var pubAddr = document.domain; // when selected: ws -> wss

var conn = new WebSocket("ws://" + addr + "/ws");

var data;

conn.onopen = function(e) {
	console.log("Connected");
	
	// Uncomment this if deploying to Heroku
	setInterval( function() {
		conn.send( "Keep connection alive!" ); // ping server every 50000 milliseconds = 50 seconds to prevent disconnection
	}, 50000 );
	
};

conn.onclose = function(e) {
	console.log("Disconnected");
};

conn.onmessage = function(e) {
	
	data = JSON.parse(e.data);
	
	switch ( data.Title ) {
		case "createPlayer": createPlayer( data.Data ); break;
		case "addPlayer": addPlayer( data.Data ); break;
		case "updatePlayer": updatePlayer( data.Data ); break;
		case "removePlayer": removePlayer( data.Data.ID ); break;
		default: console.log("error");
	}
	
};

var createPlayer = function( data ) {
	
	// Player
	playerID = data.ID + "";
	playerData = data.Pos;
	
	player = new THREE.Mesh( cube_geometry, cube_material );
	
	player.position.x = data.Pos.X;
	player.position.y = data.Pos.Y;
	player.position.z = data.Pos.Z;
	
	player.rotation.x = data.Pos.R_X;
	player.rotation.y = data.Pos.R_Y;
	player.rotation.z = data.Pos.R_Z;
	
	updatePlayerData();
	
	scene.add( player );
	
	controls = new THREE.PlayerControls( camera, player );
	controls.addEventListener( 'change', render );
	
};

var updatePlayerData = function() {
	
	playerData.X = player.position.x;
	playerData.Y = player.position.y;
	playerData.Z = player.position.z;
	
	playerData.R_X = player.rotation.x;
	playerData.R_Y = player.rotation.y;
	playerData.R_Z = player.rotation.z;
	
	conn.send(JSON.stringify({
    	Title: "updatePlayer",
    	Pos: playerData
    }));
	
}

var addPlayer = function( data ) {
	
	var cube = new THREE.Mesh( cube_geometry, cube_material );
	
	cube.position.x = data.Pos.X;
	cube.position.y = data.Pos.Y;
	cube.position.z = data.Pos.Z;
	
	cube.rotation.x = data.Pos.R_X;
	cube.rotation.y = data.Pos.R_Y;
	cube.rotation.z = data.Pos.R_Z;
	
	otherPlayers[data.ID + ""] = cube;
	
	scene.add( otherPlayers[data.ID + ""] );
	
};

var updatePlayer = function( data ) {
	
	otherPlayers[data.ID + ""].position.x = data.Pos.X;
	otherPlayers[data.ID + ""].position.y = data.Pos.Y;
	otherPlayers[data.ID + ""].position.z = data.Pos.Z;
	
	otherPlayers[data.ID + ""].rotation.x = data.Pos.R_X;
	otherPlayers[data.ID + ""].rotation.y = data.Pos.R_Y;
	otherPlayers[data.ID + ""].rotation.z = data.Pos.R_Z;
	
};

var removePlayer = function( ID ) {
	
	scene.remove( otherPlayers[ID + ""] );
	delete otherPlayers[ID + ""];
	
};
