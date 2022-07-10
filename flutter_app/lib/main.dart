import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'router/app_router.dart';

void main() {
  runApp(App());
}

class App extends StatelessWidget {
  App({Key? key}) : super(key: key);

  final _appRouter = AppRouter();

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      routerDelegate: _appRouter.delegate(),
      routeInformationParser: _appRouter.defaultRouteParser(),
      theme: ThemeData(
        fontFamily: 'DM Sans', 
        textTheme: Theme.of(context).textTheme.apply(
          bodyColor: SchejColors.black,
          displayColor: SchejColors.black,
        ),
      ),
      color: SchejColors.darkGreen,
    );
  }
}