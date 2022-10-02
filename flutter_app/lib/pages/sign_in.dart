import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_signin_button/button_list.dart';
import 'package:flutter_signin_button/button_view.dart';
import 'package:flutter_svg/svg.dart';
import 'package:provider/provider.dart';

class SignInPage extends StatefulWidget {
  final VoidCallback onSignIn;

  const SignInPage({
    Key? key,
    required this.onSignIn,
  }) : super(key: key);

  @override
  State<SignInPage> createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          children: [
            Expanded(
              child: SvgPicture.asset(
                'assets/img/schej_logo_big.svg',
                semanticsLabel: 'Schej logo',
                width: 150,
              ),
            ),
            Expanded(
              child: CustomPaint(
                painter: _CirclePainter(),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
                      padding: const EdgeInsets.only(bottom: 50),
                      child: Text(
                        'Scheduling made simple',
                        style: SchejFonts.header
                            .copyWith(color: SchejColors.white),
                      ),
                    ),
                    SignInButton(
                      Buttons.Google,
                      onPressed: _signIn,
                      text: 'Continue with Google',
                      padding: const EdgeInsets.all(5),
                      shape: RoundedRectangleBorder(
                        borderRadius: SchejConstants.borderRadius,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _signIn() async {
    final authService = context.read<AuthService>();
    bool success = await authService.signIn();
    if (success) {
      widget.onSignIn();
    } else {}
  }
}

class _CirclePainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = SchejColors.green
      ..style = PaintingStyle.fill;

    canvas.drawCircle(
        Offset(0.5 * size.width, .8 * size.height), 1.3 * size.width, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
  }
}
