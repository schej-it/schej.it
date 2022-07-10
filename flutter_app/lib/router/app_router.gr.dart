// **************************************************************************
// AutoRouteGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND

// **************************************************************************
// AutoRouteGenerator
// **************************************************************************
//
// ignore_for_file: type=lint

part of 'app_router.dart';

class _$AppRouter extends RootStackRouter {
  _$AppRouter([GlobalKey<NavigatorState>? navigatorKey]) : super(navigatorKey);

  @override
  final Map<String, PageFactory> pagesMap = {
    SignInPageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const SignInPage());
    },
    HomePageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const HomePage());
    },
    MySchejTabRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const MySchejTab());
    },
    EventsTabRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const EventsTab());
    },
    FriendsTabRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const FriendsTab());
    },
    ProfileTabRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const ProfileTab());
    },
    MySchejPageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const MySchejPage());
    },
    EventsPageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const EventsPage());
    },
    EventPageRoute.name: (routeData) {
      final pathParams = routeData.inheritedPathParams;
      final args = routeData.argsAs<EventPageRouteArgs>(
          orElse: () =>
              EventPageRouteArgs(eventId: pathParams.getString('id')));
      return MaterialPageX<dynamic>(
          routeData: routeData,
          child: EventPage(key: args.key, eventId: args.eventId));
    },
    FriendsPageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const FriendsPage());
    },
    ProfilePageRoute.name: (routeData) {
      return MaterialPageX<dynamic>(
          routeData: routeData, child: const ProfilePage());
    }
  };

  @override
  List<RouteConfig> get routes => [
        RouteConfig(SignInPageRoute.name, path: '/sign-in'),
        RouteConfig(HomePageRoute.name, path: '/', children: [
          RouteConfig('#redirect',
              path: '',
              parent: HomePageRoute.name,
              redirectTo: 'my-schej',
              fullMatch: true),
          RouteConfig(MySchejTabRoute.name,
              path: 'my-schej',
              parent: HomePageRoute.name,
              children: [
                RouteConfig(MySchejPageRoute.name,
                    path: '', parent: MySchejTabRoute.name)
              ]),
          RouteConfig(EventsTabRoute.name,
              path: 'events',
              parent: HomePageRoute.name,
              children: [
                RouteConfig(EventsPageRoute.name,
                    path: '', parent: EventsTabRoute.name),
                RouteConfig(EventPageRoute.name,
                    path: ':id', parent: EventsTabRoute.name)
              ]),
          RouteConfig(FriendsTabRoute.name,
              path: 'friends',
              parent: HomePageRoute.name,
              children: [
                RouteConfig(FriendsPageRoute.name,
                    path: '', parent: FriendsTabRoute.name)
              ]),
          RouteConfig(ProfileTabRoute.name,
              path: 'profile',
              parent: HomePageRoute.name,
              children: [
                RouteConfig(ProfilePageRoute.name,
                    path: '', parent: ProfileTabRoute.name)
              ])
        ])
      ];
}

/// generated route for
/// [SignInPage]
class SignInPageRoute extends PageRouteInfo<void> {
  const SignInPageRoute() : super(SignInPageRoute.name, path: '/sign-in');

  static const String name = 'SignInPageRoute';
}

/// generated route for
/// [HomePage]
class HomePageRoute extends PageRouteInfo<void> {
  const HomePageRoute({List<PageRouteInfo>? children})
      : super(HomePageRoute.name, path: '/', initialChildren: children);

  static const String name = 'HomePageRoute';
}

/// generated route for
/// [MySchejTab]
class MySchejTabRoute extends PageRouteInfo<void> {
  const MySchejTabRoute({List<PageRouteInfo>? children})
      : super(MySchejTabRoute.name,
            path: 'my-schej', initialChildren: children);

  static const String name = 'MySchejTabRoute';
}

/// generated route for
/// [EventsTab]
class EventsTabRoute extends PageRouteInfo<void> {
  const EventsTabRoute({List<PageRouteInfo>? children})
      : super(EventsTabRoute.name, path: 'events', initialChildren: children);

  static const String name = 'EventsTabRoute';
}

/// generated route for
/// [FriendsTab]
class FriendsTabRoute extends PageRouteInfo<void> {
  const FriendsTabRoute({List<PageRouteInfo>? children})
      : super(FriendsTabRoute.name, path: 'friends', initialChildren: children);

  static const String name = 'FriendsTabRoute';
}

/// generated route for
/// [ProfileTab]
class ProfileTabRoute extends PageRouteInfo<void> {
  const ProfileTabRoute({List<PageRouteInfo>? children})
      : super(ProfileTabRoute.name, path: 'profile', initialChildren: children);

  static const String name = 'ProfileTabRoute';
}

/// generated route for
/// [MySchejPage]
class MySchejPageRoute extends PageRouteInfo<void> {
  const MySchejPageRoute() : super(MySchejPageRoute.name, path: '');

  static const String name = 'MySchejPageRoute';
}

/// generated route for
/// [EventsPage]
class EventsPageRoute extends PageRouteInfo<void> {
  const EventsPageRoute() : super(EventsPageRoute.name, path: '');

  static const String name = 'EventsPageRoute';
}

/// generated route for
/// [EventPage]
class EventPageRoute extends PageRouteInfo<EventPageRouteArgs> {
  EventPageRoute({Key? key, required String eventId})
      : super(EventPageRoute.name,
            path: ':id',
            args: EventPageRouteArgs(key: key, eventId: eventId),
            rawPathParams: {'id': eventId});

  static const String name = 'EventPageRoute';
}

class EventPageRouteArgs {
  const EventPageRouteArgs({this.key, required this.eventId});

  final Key? key;

  final String eventId;

  @override
  String toString() {
    return 'EventPageRouteArgs{key: $key, eventId: $eventId}';
  }
}

/// generated route for
/// [FriendsPage]
class FriendsPageRoute extends PageRouteInfo<void> {
  const FriendsPageRoute() : super(FriendsPageRoute.name, path: '');

  static const String name = 'FriendsPageRoute';
}

/// generated route for
/// [ProfilePage]
class ProfilePageRoute extends PageRouteInfo<void> {
  const ProfilePageRoute() : super(ProfilePageRoute.name, path: '');

  static const String name = 'ProfilePageRoute';
}
