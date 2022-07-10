import 'package:auto_route/auto_route.dart';
import 'package:flutter/material.dart';

import '../constants/colors.dart';

class SchejAppBar extends AppBar {
  SchejAppBar({
    Key? key,
    required Widget title,
    bool isRoot = false,
  }) : super(
    key: key, 
    leading: isRoot ? null : const AutoLeadingButton(),
    title: title,
    centerTitle: false,
    foregroundColor: SchejColors.black,
    backgroundColor: SchejColors.white,
    elevation: 1,
  );
}