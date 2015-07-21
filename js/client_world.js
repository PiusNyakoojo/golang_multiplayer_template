var container, scene, camera, renderer;
	
var controls;

var sphere;

var cube_geometry = new THREE.BoxGeometry( 1, 1, 1 );
var cube_material = new THREE.MeshBasicMaterial( {color: 0x7777ff, wireframe: false} );

var player, playerID, playerData;

var otherPlayers = {};

init();
animate();

function init() {

	// Setup
	container = document.getElementById( 'container' );

	scene = new THREE.Scene();

	camera = new THREE.PerspectiveCamera( 50, window.innerWidth / window.innerHeight, 1, 1000 );
	camera.position.z = 3;

	renderer = new THREE.WebGLRenderer( { alpha: true} );
	renderer.setSize( window.innerWidth, window.innerHeight);
	// Add Objects To the Scene HERE
	
	// Sphere
	var sphere_geometry = new THREE.SphereGeometry( 1 );
	var sphere_material = new THREE.MeshNormalMaterial();
	sphere = new THREE.Mesh( sphere_geometry, sphere_material );

	sphere.position.x = 4;
	scene.add( sphere );

	// Events
	window.addEventListener( 'resize', onWindowResize, false );

	// Final touches
	container.appendChild( renderer.domElement );
	document.body.appendChild( container );
}

function animate() {
	requestAnimationFrame( animate );
	if ( controls )
		controls.update();
	render();
}

function render() {
	// Render Scene
	renderer.clear();
	renderer.render( scene , camera );
}

function onWindowResize() {

	camera.aspect = window.innerWidth / window.innerHeight;
	camera.updateProjectionMatrix();
	
	renderer.setSize( window.innerWidth, window.innerHeight );

}