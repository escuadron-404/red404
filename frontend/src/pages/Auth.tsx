import { useState } from "react";
import GoogleIcon from "../assets/icons/GoogleIcon";
import MetaIcon from "../assets/icons/MetaIcon";
import Button from "../components/AuthComponents/Button";
import Input from "../components/AuthComponents/Input";

function Auth() {
	const [requestCode, setRequestCode] = useState(false);

	function getVerificationCode(e: React.MouseEvent<HTMLButtonElement>) {
		e.preventDefault();
		setRequestCode(true);
		console.log("requesting verification codes");
	}
	function googleAuth(e: React.MouseEvent<HTMLButtonElement>) {
		e.preventDefault();
		console.log("logging with google");
	}
	function metaAuth(e: React.MouseEvent<HTMLButtonElement>) {
		e.preventDefault();
		console.log("logging with meta");
	}
	return (
		<main className="w-full h-screen flex flex-col items-center justify-center bg-base">
			<h1 className="text-3xl text-primary mb-10">red404</h1>
			<div className="auth-providers flex flex-col gap-5 mb-5 ">
				<Button type="button" provider="Google" onClick={googleAuth}>
					<GoogleIcon />
				</Button>
				<Button type="button" provider="Meta" onClick={metaAuth}>
					<MetaIcon />
				</Button>
			</div>
			{!requestCode ? (
				<form
					className="flex flex-col items-center"
					action="/backend/login"
					method="POST"
				>
					<div className="container flex flex-col gap-3">
						<label className="self-start" htmlFor="email">
							Email
						</label>
						<Input
							name="email"
							id="email"
							placeholder="Enter your email"
							type="email"
						/>
						<Button
							type="submit"
							text="Continue"
							onClick={getVerificationCode}
						/>
					</div>
				</form>
			) : (
				<form className="" action="/verify-code">
					<div className="container flex flex-col gap-3">
						<label className="self-start" htmlFor="verification-code">
							Verification code
						</label>
						<Input
							name="verification-code"
							id="verification-code"
							placeholder="Enter code"
							type="text"
						/>
						<span className="text-text-light text-sm">
							Sent you a verification code at your email
						</span>
						<Button type="button" text="Send" />
					</div>
				</form>
			)}
		</main>
	);
}
export default Auth;
