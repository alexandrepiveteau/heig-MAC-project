# Use cases

## General

1. U: `/start`
2. B: Shows a message "Hi there", displays inline keyboard with all commands like:
    * add an attempt
    * add a new route
    * find a route
    * follow someone
    * challenge someone you follow
    * see someone's profile

## Someone wants to climb a route

1. U: clicks on "add an attempt"
2. B: Adding a new attempt to an existing route.
3. B: In which gym are you climbing?
4. U: Cube
5. B: Autocompletes to show all gym matching this name
6. U: clicks on the correct gym
7. B: What is the route's name?
8. U: Autocompletes to show all matching routes in that gym
9. U: clicks on the correct route
10. B: Have you, `<Flashed, Succeeded, Failed>`?
11. U: clicks on the correct input
14. B: On a scale of 0 to 10, how would you rate the route? (custom keyboard)
15. U: clicks
16. B: What do you think this route should be graded? (inline keyboard)
17. U: clicks
18. B: Do you want to upload a video or photo of the attempt?
19. U: no or sends pic/vid
20. B: Long live the swollen forearms! :emoji:

## Someone wants to add a route

1. U: clicks on "add a route"
2. B: Adding a new route
3. B: In which gym do you want to add the route?
4. U: Cube
5. B: Autocompletes to show all gym matching this name
6. U: clicks on the correct gym
7. B: What is the route's name?
8. U: 14
9. B: What are the route's grade? (custom keyboard)
10. U: clicks
11. B: What is the colour of the holds?
12. U: purple
13. B: When was it set?
14. U: enters date (TBD)
15. B: Do you want to add a picture of the route ? (no or image)
16. U: sends image
17. B: Long live the swollen forearms! :emoji:

## Someone wants to find a route

1. U: clicks on "find a route"
2. B: Searching for routes
3. B: In which gym do you want to find the route?
4. U: Cube
5. B: Autocompletes to show all gym matching this name
6. U: clicks on the correct gym
8. B: Which grade is the route ? (Shows autocomplete and NA/IDK)
9. U: clicks
10. B: Which color is the route ? (Shows autocomplete and NA/IDK)
11. U: clicks
12. B: Shows autocomplete list
13. U: clicks
14. B: shows route profile

## Follow someone

1. U: clicks on "follow"
2. B: Following someone
3. B: What is the @username the person you want to follow? Make sure he already contacted me at least once.
4. U: @personUsername
5. Two cases:
   - B: Great! You're now following @personUsername.\
   - B: @personUsername never contacted me. Send him this link to get him started: t.me/climbot

## See someone's profile

1. U: clicks on "search user profile"
2. B: Searching user profile
3. B: What is the @username of the person you want to check out? Make sure he already contacted me at least once.
4. U: @personUsername
5. Two cases:
   - B: Here is @personUsername's profile:
      - Favourite gym: Le cube
      - Best route climbed: La mer noire (8C)
      - Followers: 868
      - Following: 37
   - B: @personUsername never contacted me. Send him this link to get him started: t.me/climbot

## Challenge someone you follow

1. U: Clicks on "challenge someone"
2. B: Challing someone you follow
3. B: Here are the people you follow: Who do you want to challenge? (If you don't see them in this list, enter their @personUsername) (Custom keyboard containing people the user follows)
4. U: Clicks @personUsername / (Enter @personUsername)
5. Two cases:
   - B: In which gym is the route you want to challenge @personUsername to? (Custom keyboard containing gyms)
   - B: @personUsername never contacted me. Send him this link to get him started: t.me/climbot
6. U: Clicks on the correct gym
7. B: What is the route's name? (If you don't see it in this list, enter it's name manually.) (Custom keyboard containing route names (sorted by recent dates?)
8. U: Clicks on the correct route
9. Two cases:
   - B: I don't recognize this route/gym combination. Try checking for spelling mistaked, or add the route to my database using /addRoute.
   - B: Great! You challenged @personUsername to climb routeName in GymName! 
10. B -> @personUsernam: @User challenged you to beat the route "routeName" in the gym "gymName". Don't forget to mark the route as climb using /climbRoute once you succeed!
