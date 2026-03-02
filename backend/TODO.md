# TODO Backend
### Multi Login
- [x] Set userToken and clientToken on register
- [x] Set clientToken and userToken on login
- [x] If logged in from other client, only update clientToken on login
- [x] Update ValidateToken to follow new pattern
- [x] Move expiry date to client token
- [x] Update RemoveToken to follow new pattern

### Password reset
- [x] Generate JWT token and verify it to get data
- [x] Send email with link to reset at
- [x] From on website sending data to endpoint
- [ ] Handle data in endpoint
- [ ] Set new password in DB
- [ ] Sign out all devices
- [ ] `(all_good ? redirect to success.html : redirect to error.html)`
