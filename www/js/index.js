// https://auth0.com/docs/client-platforms/jquery

var lock = null;
$(document).ready(function() {
   lock = new Auth0Lock('YOUR_CLIENT_ID', 'YOUR_NAMESPACE');
});


var userProfile;

$('.btn-login').click(function(e) {
  e.preventDefault();
  lock.show(function(err, profile, token) {
    if (err) {
      // Error callback
      alert('There was an error');
    } else {
      // Success callback

      // Save the JWT token.
      localStorage.setItem('userToken', token);

      // Save the profile
      userProfile = profile;
    }
  });
});


<!-- ... -->
<input type="submit" class="btn-login" />
<!-- ... -->


$.ajaxSetup({
  'beforeSend': function(xhr) {
    if (localStorage.getItem('userToken')) {
      xhr.setRequestHeader('Authorization',
            'Bearer ' + localStorage.getItem('userToken'));
    }
  }
});



$('.nick').text(userProfile.nickname);

<p>His name is <span class="nick"></span></p>


    function supportsHTML5Storage() {
        try {
            return 'localStorage' in window && window['localStorage'] !== null;
        } catch (e) {
            return false;
        }
    }



localStorage.removeItem('token');
userProfile = null;
window.location.href = "/";
