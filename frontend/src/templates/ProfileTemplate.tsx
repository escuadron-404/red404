import postData from "./PostTemplate";

const postData1 = { ...postData, postMedia: "/testPostImage2.jpg" };
const postData2 = { ...postData, postMedia: "/testPostImage3.jpg" };
const profileData = {
	userData: {
		userName: "pablopicasox",
		fullName: "Pablo Picaso",
		followers: 15254,
		following: 121,
		posts: 25,
	},
	posts: [postData1, postData, postData2],
	highlightsData: [
		{
			id: 2,
			link: "https://red404.app/highlights/321354",
			title: "beach",
		},
		{
			id: 1,
			link: "https://red404.app/highlights/321354",
			title: "paris",
		},

		{ id: 5, link: "https://red404.app/highlights/321354", title: "frankfurt" },
	],
};

export default profileData;
