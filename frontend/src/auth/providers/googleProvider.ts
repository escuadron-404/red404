type GoogleAuthMode = "login" | "register";

function googleAuth(mode: GoogleAuthMode) {
  console.log("logging with google");

  if (mode === "login") {
    window.location.assign("/home");
  }

  if (mode === "register") {
    // TODO: create a sub-page to setup the user account
    window.location.assign("/register?setup=true");
  }
}

export default googleAuth;
