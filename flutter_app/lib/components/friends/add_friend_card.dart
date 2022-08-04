import 'package:flutter/material.dart';
import 'package:flutter_app/components/user_avatar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class AddFriendCard extends StatelessWidget {
  final String name;
  final String picture;
  final bool requestAlreadySent;

  const AddFriendCard({
    Key? key,
    required this.name,
    this.picture = 'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
    this.requestAlreadySent = false,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: ListTile(
        leading: UserAvatar(src: picture),
        title: Text(name, style: SchejFonts.subtitle),
        trailing: IconButton(
          onPressed: () {},
          icon: requestAlreadySent
              ? const Icon(MdiIcons.accountArrowRight, color: SchejColors.gray)
              : const Icon(MdiIcons.accountPlus, color: SchejColors.black),
        ),
      ),
    );
  }
}
