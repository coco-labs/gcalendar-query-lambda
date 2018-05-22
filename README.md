# GCalendar Lambda Query

This lambda function processes Google Calendar events and counts how many of them have past.

For example, you might pass a query: "Test event" and it will count how many of these events are present in the calendar specified.

To deploy this you will need to get an access token from the oauth2 playground: https://developers.google.com/oauthplayground/

Make sure you select the Google Calendar API V3 from the dropdown.

Then replace the access_token for the `GCALENDAR_API_KEY` in the environment variable this lambda expects.
You will also need a calendar id, this normally is the email associated to the calendar account.

Run `make` to build the golang executable. And then `serverless deploy` to deploy to your AWS Account.

You will get an url where you can start sending requests through Postman.