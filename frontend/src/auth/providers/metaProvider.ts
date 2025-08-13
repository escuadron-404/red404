function metaAuth(e: React.MouseEvent<HTMLButtonElement>) {
	e.preventDefault();
	console.log("logging with meta");
	// TODO: refactor this, for consistency we should be able to use hooks here
	window.location.assign("/home");
}

export default metaAuth;
