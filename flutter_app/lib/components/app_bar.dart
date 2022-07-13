import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/constants/fonts.dart';

import '../constants/colors.dart';

class SchejAppBar extends AppBar {
  SchejAppBar({
    Key? key,
    String titleString = '',
    Widget? title,
    List<Widget>? actions,
    bool underline = true,
    bool isRoot = false,
  }) : super(
    key: key,
    leading: isRoot ? null : const AutoLeadingButton(),
    title: titleString.isNotEmpty
        ? Text(titleString, style: SchejFonts.header)
        : title,
    centerTitle: false,
    actions: actions,
    foregroundColor: SchejColors.black,
    backgroundColor: SchejColors.white,
    elevation: underline ? 1 : 0,
  );
}
