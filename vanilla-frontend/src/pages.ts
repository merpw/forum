export const LoginSignup = (): string => {
  return `
    <section class="login">	
	<div class="container">
	<div class="forms">
		<div class="form login">
			<span class="title">Login</span>
			<form (submit)="onSubmit()" id="login-form">
				<div class="input-field">
					<input id="username-login" type="text" placeholder="Nickname/Email" required autofocus>
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
						<label for="forgot-pw" class="text">Remember me</label>
					</div>
					<a href="https://www.youtube.com/watch?v=dQw4w9WgXcQ" class="text" id="forgot-pw">Forgot password?</a>
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

		<div class="form signup" style="display: none;">
			<span class="title">Registration</span>

			<form (submit)="onSubmit()" id="signup-form">
        <label for="username-register"></label>
				<div class="input-field">
					<input id="username-register" type="text" placeholder="Nickname" required autofocus>
					<i class="uil uil-user-circle"></i>
				</div>

				<div class="input-field">
          <label for="email"></label>
					<input id="email" type="email" placeholder="Email" required>
					<i class="uil uil-envelope icon"></i>
				</div>
				<div class="name-fields">

					<div class="input-field">
            <label for="first-name"></label>
						<input id="first-name" type="text" placeholder="First name" required>
						<i class="uil uil-user"></i>
					</div>
		  
					<div class="input-field">
            <label for="last-name"></label>
						<input id="last-name" type="text" placeholder="Last name" required>
						<i class="uil uil-user"></i>
					</div>

			  	</div>

				<div class="input-field">
				<label for="age"></label>
				<input type="date" id="age" name="age"
       			value="Age"
       			min="1900-01-01" max="2020-01-01" required>
					<i class="uil uil-calender"></i>
				</div>

				<div class="input-field">
				<label for="gender"></label>
					<input list="genders" name="gender" id="gender" placeholder="Gender" required>
            <datalist id="genders">
              <option value="Male">
              <option value="Female">
              <option value="Other">
            </datalist>
					<i class="uil uil-android-alt"></i>
				</div>

				<div class="input-field">
          <label for="password-register"></label>
					<input id="password-register" type="password" class="password" placeholder="Password" required>
					<i class="uil uil-lock icon"></i>
				</div>
				
				<div class="input-field">
          <label for="password-register-repeat"></label>
					<input id="password-register-repeat" type="password" class="password" placeholder="Repeat Password" required>
					<i class="uil uil-lock icon"></i>
					<i class="uil uil-eye-slash showHidePw"></i>
				</div>

				<div class="input-field button">
          <label for="submit"></label>
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
        <a id="topnav-chat"><i class='bx bx-chat' ></i> Chat</a>
		<a id="topnav-post" href="#post">
		<i class='bx bx-duplicate' ></i> Post
		</a>
        <a id="topnav-home" href="#index">
        <i class='bx bx-home'></i> Home
        </a>
        <a id="topnav-logout" class="logout" href="#logout"><i class='bx bx-log-out'></i> Logout</a>
    </div>

    <div id="chatlist" class="chatlist">
        <a href="#close-chat" id="chat-close" class="closebtn"><i class='bx bx-x'></i></a>
        <h2 id="online-title" style="color: black; margin-left: 20px;">Online: </h2>
        <ul id="online-users" class="chat-users"></ul>
        <h2 id="offline-title" style="color: black; margin: 10px 0px 0px 20px;">Offline: </h2>
        <ul id="offline-users" class="chat-users"></ul>
    </div>


<main id="main">
<div id="chat-area"></div>

  <section id="feed">
  <section id="create-post" class="close"></section>
      <div class="categories">
      <span class="title">Categories</span>
        <div class="category-selection">
          <h3 class="category-title" id="category-facts">#facts</h3>
          <h3 class="category-title" id="category-rumors">#rumors</h3>
          <h3 class="category-title" id="category-other">#other</h3>
        </div>
      </div>
      <div id="posts-display"></div>
  </section>
</main>
`
}

export const postForm = (): string => {
  return `
<span class="title">New post</span>
<form id="post-form" (submit)="onSubmit()">
	<input id="post-title" placeholder="Title"></input>
	<textarea id="post-content" 
            rows=5 
            style="resize: none;" 
            type="text" 
            placeholder="Write something..."></textarea>
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
		<textarea id="comment-content" 
              rows=5 maxlength="250" 
              style="resize: none;" 
              type="text"
              placeholder="Write something..."></textarea>
		<div id="create-post-footer">
			<input class="comment-submit" type="submit" value="Comment">
		</div>	
	</form>
	`
}

export const errorPage = (errorcode: number): string => {
  let errorHTML: string
  switch (errorcode) {
    case 400:
      errorHTML = `
			<h1 class="error-code">${errorcode.toString()}:</h1>
			<br>
			<h1 class="error-title">Bad request</h1>
		`
      break
    case 401:
      errorHTML = `
			<h1 class="error-code">${errorcode.toString()}:</h1>
			<br>
			<h1 class="error-title">Unauthorized.</h1>
		`
      break
    case 404:
      errorHTML = `
			<h1 class="error-code">${errorcode.toString()}:</h1>
			<br>
			<h1 class="error-title">Not found.</h1>
		`
      break
    case 405:
      errorHTML = `
			<h1 class="error-code">${errorcode.toString()}:</h1>
			<br>
			<h1 class="error-title">Method not allowed.</h1>
		`
      break
    default:
      errorHTML = `
			<h1 class="error-code">${errorcode.toString()}:</h1>
			<br>
			<h1 class="error-title">Internal server error.</h1>
		`
  }
  return errorHTML
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
