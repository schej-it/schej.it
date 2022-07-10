import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/constants/fonts.dart';

import '../constants/colors.dart';

class SchejAppBar extends AppBar {
  SchejAppBar({
    Key? key,
    required String title,
    bool isRoot = false,
  }) : super(
    key: key, 
    leading: isRoot ? null : const AutoLeadingButton(),
    title: Text(title, style: SchejFonts.header),
    centerTitle: false,
    foregroundColor: SchejColors.black,
    backgroundColor: SchejColors.white,
    elevation: 1,
  );
}