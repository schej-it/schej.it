import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class CompareSchejCard extends StatelessWidget {
  final String name;
  final String pic;
  final bool added;
  final void Function(bool) onToggle;

  const CompareSchejCard({
    Key? key,
    required this.name,
    this.pic = 'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
    this.added = false,
    required this.onToggle,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: const Duration(milliseconds: 100),
      decoration: SchejConstants.listTileDecoration
          .copyWith(color: added ? SchejColors.green : SchejColors.white),
      child: ListTile(
        // tileColor: added ? SchejColors.green : SchejColors.white,
        leading: CircleAvatar(
          backgroundImage: NetworkImage(pic),
        ),
        title: Text(name,
            style: SchejFonts.subtitle.copyWith(
                color: added ? SchejColors.white : SchejColors.pureBlack)),
        trailing: IconButton(
          onPressed: () {
            onToggle(!added);
          },
          icon: added
              ? const Icon(MdiIcons.close, color: SchejColors.white)
              : const Icon(MdiIcons.plus, color: SchejColors.black),
        ),
      ),
    );
  }
}
