import 'package:auto_route/auto_route.dart';
import 'package:flutter_app/pages/friends/compare_schej.dart';
import 'package:flutter_app/router/auth_guard.dart';
import '../pages/events/event.dart';
import '../pages/events/events.dart';
import '../pages/events/events_tab.dart';
import '../pages/friends/friends.dart';
import '../pages/friends/friends_tab.dart';
import '../pages/my_schej/my_schej.dart';
import '../pages/my_schej/my_schej_tab.dart';
import '../pages/profile/profile.dart';
import '../pages/profile/profile_tab.dart';
import '../pages/sign_in.dart';
import '../pages/home.dart';

@MaterialAutoRouter(
  routes: <AutoRoute>[
    AutoRoute(path: '/sign-in', page: SignInPage),
    AutoRoute(
      path: '/',
      page: HomePage,
      guards: [AuthGuard],
      children: [
        AutoRoute(
          path: 'my-schej',
          page: MySchejTab,
          initial: true,
          children: [
            AutoRoute(path: '', page: MySchejPage),
          ],
        ),
        AutoRoute(
          path: 'events',
          page: EventsTab,
          children: [
            AutoRoute(path: '', page: EventsPage),
            AutoRoute(path: ':id', page: EventPage),
          ],
        ),
        AutoRoute(
          path: 'friends',
          page: FriendsTab,
          children: [
            AutoRoute(path: '', page: FriendsPage),
            AutoRoute(path: 'schej', page: CompareSchejPage),
          ],
        ),
        AutoRoute(
          path: 'profile',
          page: ProfileTab,
          children: [
            AutoRoute(path: '', page: ProfilePage),
          ],
        ),
      ],
    ),
  ],
)
class $AppRouter {}
