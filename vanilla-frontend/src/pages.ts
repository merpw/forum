export const LoginSignup = (): string => {
  return `
    <section class="login">	
	<div class="container">
	<div class="forms">
		<div class="form login">
			<span class="title">Login</span>
			<form (submit)="onSubmit()" id="login-form">
				<div class="input-field">
					<input id="username-login" type="text" placeholder="Nickname/Email" required>
					<i class="uil uil-user"></i>
				</div>

				<div class="input-field">
					<input id="password-login" type="password" class="password" placeholder="Password" required>
					<i class="uil uil-lock icon"></i>
					<i class="uil uil-eye-slash showHidePw"></i>
				</div>

				<div class="checkbox-text">
					<div class="checkbox-content">
						<input type="checkbox" id="logCheck">
						<label for="logCheck" class="text">Remember me</label>
					</div>
					<a href="https://www.youtube.com/watch?v=dQw4w9WgXcQ" class="text">Forgot password?</a>
				</div>

				<div class="input-field button">
					<input type="submit" value="Login Now" class="login-btn">
				</div>

				<div class="login-signup">
					<span class="text">Not registered?
						<a href="#" class="text signup-link">Sign up now!</a>
					</span>
				</div>
			</form>
		</div>

		<div class="form signup">
			<span class="title">Registration</span>

			<form (submit)="onSubmit()" id="signup-form">
				<div class="input-field">
					<input id="username-register" type="text" placeholder="Nickname" required>
					<i class="uil uil-user-circle"></i>
				</div>

				<div class="input-field">
					<input id="email" type="email" placeholder="Email" required>
					<i class="uil uil-envelope icon"></i>
				</div>

				<div class="name-fields">

					<div class="input-field">
						<input id="first-name" type="text" placeholder="First name" required>
						<i class="uil uil-user"></i>
					</div>
		  
					<div class="input-field">
						<input id="last-name" type="text" placeholder="Last name" required>
						<i class="uil uil-user"></i>
					</div>

			  	</div>

				<div class="input-field">
				
				<input type="date" id="age" name="age"
       			value="Age"
       			min="1900-01-01" max="2020-01-01" required>
					<i class="uil uil-calender"></i>
				</div>

				<div class="input-field">
					<input id="gender" type="text" placeholder="Gender" required>
					<i class="uil uil-android-alt"></i>
				</div>

				<div class="input-field">
					<input id="password-register" type="password" class="password" placeholder="Password" required>
					<i class="uil uil-lock icon"></i>
				</div>
				
				<div class="input-field">
					<input id="password-register-repeat" type="password" class="password" placeholder="Repeat Password" required>
					<i class="uil uil-lock icon"></i>
					<i class="uil uil-eye-slash showHidePw"></i>
				</div>

				<div class="input-field button">
					<input type="submit" value="Sign Up" class="signup">
				</div>

				<div class="login-signup">
					<span class="text">Already registered?
						<a href="#" class="text login-link">Log in now!</a>
					</span>
				</div>
			</form>
		</div>
	</div>
</div>
</section>
`
}

export const Index = (): string => {
  return `
    <div class="topnav" id="myTopnav">
        <a href="#chat" onclick="openNav()"><i class='bx bx-chat' ></i> Chat</a>
		<a id="topnav-post" href="#post">
		<i class='bx bx-duplicate' ></i> Post
		</a>
        <a id="topnav-home" href="#index">
        <i class='bx bx-home'></i> Home
        </a>
        <a id="topnav-logout" class="logout" href="#logout"><i class='bx bx-log-out'></i> Logout</a>
    </div>

    <div id="chatlist" class="chatlist">
        <a href="javascript:void(0)" class="closebtn" onclick="closeNav()">&times;</a>
        <h2 style="color: white; margin-left: 20px;">Chats:</h2>
        <ul class="chat-users">
            <li class="online">Online user
                <i class='bx bx-message-dots'></i>
            </li>
            <li class="offline">Offline user</li>
        </ul>
    </div>

<section id="create-post" class="close">

</section>
<section class="feed">
        <div class="categories">
		<span class="title">Categories:</span>
			<div class="category-selection">
				<h3 class="category-title" id="category-facts">#facts</h3>
				<h3 class="category-title" id="category-rumors">#rumors</h3>
				<h3 class="category-title" id="category-other">#other</h3>
			</div>
		</div>
		<div id="posts-display"></div>
</section>
`
}

export const postForm = (): string => {
  return `
<span class="title">New post:</span>
<form id="post-form" (submit)="onSubmit()">
	<div id="create-post-flexbox">
		<input id="post-title" placeholder="Title"></input>
	</div>

	<textarea id="post-content" rows=5 style="resize: none;" type="text"></textarea>

	<div id="create-post-footer">
		<select id="post-category" name="post-category">
			<option value="facts">Facts</option>
			<option value="rumors">Rumors</option>
			<option value="other">Other</option>
		</select>
		<input id="post-submit" type="submit" value="Post">
	</div>	
</form>
	`
}

export const commentForm = (postId: string): string => {
  return `
	<form id="comment-form-${postId}" (submit)="onSubmit()">	
		<textarea id="comment-content" rows=5 maxlength="250" style="resize: none;" type="text"></textarea>
		<div id="create-post-footer">
			<input class="comment-submit" type="submit" value="Post">
		</div>	
	</form>
	`
}

/*
<div class="chat-window hide">
    <ul class="messages-display"></ul>
    <form class="chat-form">
        <input name="chat-input" class="chat-input">
        <button name="send-btn" class="send-btn" type="submit"><i class='bx bx-send'></i></button>
    </form>
</div>
*/
