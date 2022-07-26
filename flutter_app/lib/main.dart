import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:flutter_app/router/auth_guard.dart';
import 'package:provider/provider.dart';

import 'constants/constants.dart';
import 'constants/fonts.dart';

void main() {
  runApp(const App());
}

class App extends StatefulWidget {
  const App({Key? key}) : super(key: key);

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  final _authService = AuthService();
  late final _appRouter = AppRouter(authGuard: AuthGuard(_authService));

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => _authService,
      child: _buildMaterialApp(),
    );
  }

  Widget _buildMaterialApp() {
    return MaterialApp.router(
      routerDelegate: _appRouter.delegate(),
      routeInformationParser: _appRouter.defaultRouteParser(),
      theme: ThemeData(
        cupertinoOverrideTheme: const CupertinoThemeData(
          primaryColor: SchejColors.darkGreen,
        ),
        primaryColor: SchejColors.darkGreen,
        fontFamily: 'DM Sans',
        textTheme: Theme.of(context).textTheme.apply(
              bodyColor: SchejColors.black,
              displayColor: SchejColors.black,
            ),
        textSelectionTheme: const TextSelectionThemeData(
          cursorColor: SchejColors.black,
          selectionColor: SchejColors.darkGreen,
          selectionHandleColor: SchejColors.darkGreen,
        ),
        inputDecorationTheme: Theme.of(context).inputDecorationTheme.copyWith(
              // isDense: true,
              contentPadding: const EdgeInsets.all(0),
              hintStyle:
                  SchejFonts.subtitle.copyWith(color: SchejColors.darkGray),
              enabledBorder: OutlineInputBorder(
                borderRadius: SchejConstants.borderRadius,
                borderSide:
                    const BorderSide(width: 1, color: SchejColors.lightGray),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: SchejConstants.borderRadius,
                borderSide:
                    const BorderSide(width: 2, color: SchejColors.darkGreen),
              ),
              prefixIconColor: MaterialStateColor.resolveWith((states) {
                if (states.contains(MaterialState.focused)) {
                  return SchejColors.darkGreen;
                }
                return SchejColors.darkGray;
              }),
            ),
      ),
      color: SchejColors.darkGreen,
    );
  }
}
