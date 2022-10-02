import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_app/router/app_router.gr.dart';
import 'package:flutter_app/router/auth_guard.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:provider/provider.dart';
import 'package:requests/requests.dart';
// ignore: implementation_imports
import 'package:requests/src/cookie.dart';

import 'constants/constants.dart';
import 'constants/fonts.dart';

void main() {
  // SystemChrome.setSystemUIOverlayStyle(const SystemUiOverlayStyle(
  //   systemNavigationBarColor: Colors.white,
  //   statusBarColor: Colors.white,
  // ));
  runApp(const App());
}

class App extends StatefulWidget {
  const App({Key? key}) : super(key: key);

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  final _apiService = ApiService();
  late final _authService = AuthService(apiService: _apiService);
  late final _appRouter = AppRouter(authGuard: AuthGuard(_authService));

  bool _initialized = false;

  @override
  void initState() {
    super.initState();

    init();
  }

  Future<void> init() async {
    // Read session cookie from storage if it exists
    const storage = FlutterSecureStorage();
    final String? sessionCookieString =
        await storage.read(key: 'sessionCookie');
    if (sessionCookieString != null) {
      final CookieJar cookieJar =
          CookieJar.parseCookiesString(sessionCookieString);
      await Requests.setStoredCookies(
        Requests.getHostname(ApiService.serverAddress),
        cookieJar,
      );
    }

    // Sign in silently then set app to initialized
    await _authService.signInSilently();
    setState(() {
      _initialized = true;
    });
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (context) => _authService,
        ),
        ChangeNotifierProvider(
          create: (context) => _apiService,
        ),
      ],
      child: _buildMaterialApp(),
    );
  }

  Widget _buildMaterialApp() {
    if (!_initialized) {
      // TODO: replace with a loading screen
      return Container(
        color: SchejColors.white,
      );
    }

    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      routerDelegate: _appRouter.delegate(),
      routeInformationParser: _appRouter.defaultRouteParser(),
      theme: ThemeData(
        cupertinoOverrideTheme: const CupertinoThemeData(
          primaryColor: SchejColors.darkGreen,
        ),
        primaryColor: SchejColors.darkGreen,
        fontFamily: 'DM Sans',
        textTheme: Theme.of(context).textTheme.apply(
              bodyColor: SchejColors.pureBlack,
              displayColor: SchejColors.pureBlack,
            ),
        textSelectionTheme: const TextSelectionThemeData(
          cursorColor: SchejColors.pureBlack,
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
                    const BorderSide(width: 1, color: SchejColors.lightGray),
              ),
              prefixIconColor: _conditionalColor(
                focused: SchejColors.darkGreen,
                unfocused: SchejColors.darkGray,
              ),
              filled: true,
              fillColor: _conditionalColor(
                focused: SchejColors.white,
                unfocused: SchejColors.offWhite,
              ),
            ),
        checkboxTheme: CheckboxThemeData(
          fillColor: MaterialStateProperty.resolveWith((states) {
            // const interactiveStates = <MaterialState>{
            //   MaterialState.pressed,
            //   MaterialState.hovered,
            //   MaterialState.focused,
            // };
            // if (states.any(interactiveStates.contains)) {
            //   return SchejColors.darkGreen;
            // }
            return SchejColors.darkGreen;
          }),
          shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(Radius.circular(3)),
          ),
        ),
      ),
      color: SchejColors.white,
    );
  }

  Color _conditionalColor({dynamic focused, dynamic unfocused}) {
    return MaterialStateColor.resolveWith((states) {
      if (states.contains(MaterialState.focused)) {
        return focused;
      }
      return unfocused;
    });
  }
}
